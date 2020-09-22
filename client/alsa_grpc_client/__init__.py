import queue
import threading
import logging
import grpc
from alsa_grpc_client.grpc_gen import alsamixer_pb2
from alsa_grpc_client.grpc_gen import alsamixer_pb2_grpc

__version__ = "0.0.8"

_LOGGER = logging.getLogger(__name__)

class AlsaControl:

    def __init__(self, card, name, volume, req_queue):
        self.card = card
        self.name = name
        self.volume = volume
        self.on_update = None
        self.req_queue = req_queue

    def subscribe(self, cb):
        self.on_update = cb

    def apply(self, volume):
        self.volume = volume
        self.on_update()

    def set_volume(self, volume_level):
        _LOGGER.debug('Setting volume for %s to %f', self.name, volume_level)
        self.req_queue.put_nowait((self.name, [int(round(volume_level * 100))]))


class AlsaClient:
    def __init__(self, host, port, on_connect):
        self.host = host
        self.port = port
        self.queue = queue.SimpleQueue()
        self.connected = threading.Event()
        self.disconnect_requested = threading.Event()
        self.on_connect = on_connect
        self.controls = {}

    def connect(self):
        if self.connected.is_set():
            _LOGGER.warning("Already connected to %s:%d", self.host, self.port)
        else:
            thread = threading.Thread(target=self.__run)
            thread.start()

    def disconnect(self):
        # Set disconnected requested flag
        self.disconnect_requested.set()

        # If connected, queue poison pill (will exit __yield_results and then break loop in __run() and exit thread)
        if self.connected.is_set():
            self.queue.put(None)
        else:
            _LOGGER.warning("Already disconnected from %s:%d", self.host, self.port)

    def __run(self):
        while True:
            _LOGGER.info("Connecting to %s:%d", self.host, self.port)
            channel = grpc.insecure_channel(self.host + ':' + str(self.port))
            stub = alsamixer_pb2_grpc.AlsamixerStub(channel)
            try:
                for response in stub.Communicate(self.__yield_requests()):

                    # Always set connected flag (though the flag will only transition once per connection)
                    self.connected.set()

                    for update in response.controls:
                        if update.name in self.controls:
                            self.controls[update.name].apply(update.volume)
                        else:
                            ctrl = AlsaControl(response.card, update.name, update.volume, self.queue)
                            self.controls[update.name] = ctrl
                            self.on_connect(ctrl)

                # Reaches here when __yield_requests() has already exited cleanly (occurs when disconnect() was invoked)
                break
            except grpc.RpcError as e:
                if e.code() == grpc.StatusCode.UNAVAILABLE:

                    # Clear connected flag
                    self.connected.clear()

                    # Queue poison pill in order to cleanly exit generator created by __yield_requests()
                    self.queue.put(None)
                else:
                    _LOGGER.exception("GRPC error from %s:%d", self.host, self.port)

            if self.disconnect_requested.wait(5):
                break
            else:
                _LOGGER.info("Attempting to reconnect to %s:%d", self.host, self.port)
        self.connected.clear()
        _LOGGER.info("Exiting run thread for %s:%d", self.host, self.port)

    def __yield_requests(self):
        while True:
            req = self.queue.get()
            if req is None:
                break
            control, volume = req
            if self.connected.is_set():
                update = alsamixer_pb2.Request(control=control, volume=volume)
                yield update
            else:
                _LOGGER.warning("Exiting set_volume request for %s (not connected to %s:%d)", control, self.host, self.port)
        _LOGGER.info("Exiting yield_requests generator for %s:%d", self.host, self.port)


if __name__ == '__main__':
    logging.basicConfig()


    def on_connect(ctrl):
        print('Connected to', ctrl.name, ctrl.volume)

        def on_update():
            print('Received update for', ctrl.name)
        ctrl.subscribe(on_update)



    client = AlsaClient("192.168.0.77", 50051, on_connect)

    client.connect()
    # sleep(1)
    # alsa_grpc_client.set_volume(10)
    # sleep(5)
    # alsa_grpc_client.disconnect()
    # sleep(5)
    # alsa_grpc_client.connect()
    # sleep(.1)
    # alsa_grpc_client.set_volume(0.2)
    # sleep(5)
    # alsa_grpc_client.disconnect()


