import signal
import socket
import logging

from common.protocol import ServerProtocol
from common.utils import store_bets

from common.utils import Bet

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._server_socket.settimeout(5.0)
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
            # TODO: Modify this program to handle signal to graceful shutdown
            # the server
            while self._is_running:
                self._client = self.__accept_new_connection()
                self.__handle_client_connection(self._client)
                self._client_socket = None
        except OSError as _:
            self._is_running = False
        finally:
            self.__shutdown()

    def __handle_client_connection(self, client):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:

            serialized_bet = client.receive_client_info()
            bets, read = self.__create_bet_from_message(serialized_bet)
            if len(bets) == 0: 
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: ${read}")
                client.send_confirmation(False)
                return

            store_bets(bets)   
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
            client.send_confirmation(True)
        except OSError as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
        finally:
            client.shutdown()
            self._client = None

    def __create_bet_from_message(self, message: str):
        bets = []
        errors = 0
        try:
            for bet in message.split('\n'):
                bet_parts = bet.split(',')
                if bet_parts == ['']:
                    continue
                if len(bet_parts) != 6:
                    logging.error(f"action: parse_bet | result: fail | error: invalid_bet_format | bet_parts: {bet_parts}")
                    errors += 1
                    continue

                bet = Bet(bet_parts[0], bet_parts[1], bet_parts[2], bet_parts[3], bet_parts[4], bet_parts[5])
                bets.append(bet)
        
            return bets, errors
        except Exception as e:
            logging.error(f"action: parse_bet | result: fail | error: {e}")
            return bets, errors

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        protocol = ServerProtocol(c)
        return protocol

    def __handle_sigterm(self, signum, frame):
        try:
            logging.info("action: shutdown | result: in_progress")
            self.__shutdown()
            logging.info("action: shutdown | result: success")
        except Exception as e:
            logging.error(f"action: shutdown | result: fail | error: {e}")
            return

    def __shutdown(self):
        self._client_socket = self.__close_skt(self._client_socket)
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
