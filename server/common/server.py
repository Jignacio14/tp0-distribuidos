import signal
import socket
import logging

from common.protocol import ServerProtocol
from common.utils import store_bets, Bet

DELIMITER = ','

class Server:
    def __init__(self, port, listen_backlog):
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._client = None
        self._is_running = True
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """
        try:
            while self._is_running:
                self._client = self.__accept_new_connection()
                self.__handle_client_connection(self._client)
                self._client = None
        except OSError as skt_err:
            logging.debug(f"action: server_loop | error: {skt_err}")
            self._is_running = False
        finally:
            self.__shutdown()

    def __handle_client_connection(self, client: ServerProtocol):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        serialized_bet = client.receive_client_info()
        bet = self.__create_bet_from_message(serialized_bet)
        if not bet: 
            logging.error("action: receive_message | result: fail | error: invalid_bet")
            client.send_confirmation(False)
            return

        store_bets([bet])   
        client.send_confirmation(True)
        logging.info(f"action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}")

    def __create_bet_from_message(self, message: str):
        msg_parts = message.split(DELIMITER)
        if len(msg_parts) != 6:
            return None
        return Bet(*msg_parts)

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return ServerProtocol(c)

    def __handle_sigterm(self, signum, frame):
        try:
            logging.info("action: shutdown | result: in_progress")
            self.__shutdown()
            logging.info("action: shutdown | result: success")
        except Exception as e:
            logging.error(f"action: shutdown | result: fail | error: {e}")
            return

    def __shutdown(self):
        if self._client:
            self._client.shutdown()
        self._server_socket = self.__close_skt(self._server_socket)
        self._is_running = False

    def __close_skt(self, skt): 
        if not skt:
            return
        try:
            skt.close()
        except Exception as e:
            logging.error(f"action: close_socket | result: fail | error: {e}")
        return
