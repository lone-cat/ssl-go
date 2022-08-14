FROM golang:latest

ARG UID
ARG GID

# add user with rights to external folder
# -o means that user(group) will be created even if it already exists with specified id
# -g - sets specific group id number
# -m - creates home directory
# -u - sets specific user id number
# -g (for user) specify group name
RUN groupadd -o -g ${GID} hostuser && useradd -o -m -u ${UID} -g hostuser hostuser

USER hostuser