#!/bin/bash
# Lanza los clientes de Go en background, cada uno con su propio log.
# Uso: ./run_clients.sh

LOG_DIR="logs"
mkdir -p "$LOG_DIR"

# Lista de clientes: nombre y argumentos
declare -a CLIENTS=(
  "client1 2 1"
  "client2 2 2"
  "client3 3 3"
  "client4 1 4"
  "client5 2 5"
  "client6 3 6"
  "client7 3 7"
  "client8 3 8"

)

PIDS=()

for entry in "${CLIENTS[@]}"; do
  read -r name arg1 arg2 <<< "$entry"
  logfile="$LOG_DIR/${name}.log"
  echo "Lanzando $name -> $logfile"
  go run cmd/main.go "$name" "$arg1" "$arg2" > "$logfile" 2>&1 &
  PIDS+=($!)
done

echo "Clientes lanzados con PIDs: ${PIDS[*]}"
echo "Logs en: $LOG_DIR/"

# Opcional: espera a que todos terminen antes de salir del script
wait