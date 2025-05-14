#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <sys/socket.h>
#include <netinet/in.h>


int main() {
  // Declare variables
  int retVal = 0;
  struct sockaddr_in bindAddr;
  int tcpSocket;
  int clientSocket = 0;
  const char *hello = "HTTP/1.1 200 OK\r\nContent-Type text/plain\r\nContent-Length 14\r\n\r\nHello, World!";

  // Initialize variables
  memset(&bindAddr, 0, sizeof(bindAddr));
  tcpSocket = socket(AF_INET, SOCK_STREAM, 0);

  if (-1 == tcpSocket) {
    perror("Failed to make socket");
    return 1;
  }
  printf("socket creation succeeded\n");

  bindAddr.sin_port = htons(8080);
  bindAddr.sin_family = AF_INET;
  bindAddr.sin_addr.s_addr = INADDR_ANY;

  if (bind(tcpSocket, (const struct sockaddr*)&bindAddr, sizeof(bindAddr))
      < 0) {
    perror("Failed to bind socket");
    retVal = 1;
    goto exit;
  }
  printf("bind succeeded\n");

  if (listen(tcpSocket, SOMAXCONN) < 0) {
    perror("Failed to listen on socket");
    retVal = 1;
    goto exit;
  }
  printf("listen succeeded\n");

  for (;;) {
    printf("waiting for connections...\n");
    clientSocket = accept(tcpSocket, NULL, NULL);
    printf("connection made!\n");

    char buffer[1024];
    memset(buffer, 0, sizeof(buffer));

    ssize_t n = recv(clientSocket, buffer, sizeof(buffer)-1, 0);
    if (n < 0) {
      perror("Failed to read from socket");
      retVal = 1;
      goto exit;
    }
    else if (n == 0) {
      printf("read everything!\n");
      break;
    }

    printf("Request:\n%s", buffer);
    send(clientSocket, hello, strlen(hello), 0);
    close(clientSocket);
    break;
  }



exit:
  close(tcpSocket);
  return retVal;
}
