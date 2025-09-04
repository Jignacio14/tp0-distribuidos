import signal
import socket
import logging

from common.protocol import ServerProtocol
from common.utils import has_won, load_bets, store_bets

from common.utils import Bet

DELIMITER = ','

class Server:
    def __init__(self, port, listen_backlog, clients_num):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._client : ServerProtocol = None
        self._is_running = True
        self._clients = {}
        self._total_clients = clients_num
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
                if len(self._clients) == self._total_clients:
                    logging.info("action: sorteo | result: success")
                    self.__send_winners_to_all_clients()
                    self.__shutdown_clients()
                    self._is_running = False
                    return
        except OSError as _:
            self._is_running = False
        finally:
            self.__shutdown()

    def __handle_client_connection(self, client: ServerProtocol):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        current_client_id = client.get_agency_id()
        
        if current_client_id == -1:
            client.shutdown()
            return
        
        self._clients[current_client_id] = self._clients.get(current_client_id, client)
        self.__process_client_bets(client)

    def __process_client_bets(self, client: ServerProtocol):
        receiving_bets = True

        while receiving_bets:
            keep_reading, msg = client.receive_batch()
            receiving_bets = keep_reading
            if not keep_reading:
                break
            if msg == '':
                break
            bets, errors = self.__create_bet_from_message(msg)
            if errors > 0: 
                logging.error(f"action: apuesta_recibida | result: fail | cantidad: {errors}")
                client.send_bad_bets(errors)
                return
            store_bets(bets)   
            logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")
            client.send_batches_received_successfully(len(bets))

    def __inform_winners(self, agency_id: str) -> list[Bet]:
        return [bet.document for bet in load_bets() if bet.agency == agency_id and has_won(bet)]

    def __create_bet_from_message(self, message: str):
        bets = []
        errors = 0
        for bet_str in message:
            splited = bet_str.split(DELIMITER)
            if len(splited) != 6:
                logging.error(f"action: parse_bet | result: fail | error: invalid_bet_format | bet_parts: {splited}")
                errors += 1
                continue
            bets.append(Bet(*splited))
        return bets, errors

    def __send_winners_to_all_clients(self):
        for client_id in self._clients.keys():
            winners = self.__inform_winners(client_id)
            try:
                self._clients[client_id].send_winners(winners)
                logging.info(f"action: informar_ganadores | result: success | cantidad: {len(winners)}")
            except OSError as e:
                logging.error(f"action: informar_ganadores | result: fail | error: {e}")

    def __shutdown_clients(self):
        for client in self._clients.values():
            client.shutdown()
        self._clients = {}

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
        self._server_socket = self.__close_skt(self._server_socket)
        self.__shutdown_clients()
        self._is_running = False

    def __close_skt(self, skt): 
        if not skt:
            return
        try:
            skt.close()
        except Exception as e:
            logging.error(f"action: close_socket | result: fail | error: {e}")
        return
