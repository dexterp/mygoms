# kick:render
#!/bin/bash

# start name service caching
#echo "Starting nscd"
#/usr/sbin/nscd
#if [ $$? -ne 0 ]; then
#  echo "Failed to start nscd: $$status"
#  exit $$status
#fi

# start daemon
echo "Starting ${PROJECT_NAME}server"
exec /usr/local/bin/${PROJECT_NAME}server -l ":9090"
if [ $$? -ne 0 ]; then
  echo "Failed to start ${PROJECT_NAME}server: $$status"
  exit $$status
fi