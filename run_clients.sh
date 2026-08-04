
#!/bin/bash
# Lanza los 5 clientes de Go, cada uno en su propio panel de tmux.
# Uso: ./run_clients.sh
 
SESSION="clients"
 
# Si ya existe una sesion con ese nombre, la matamos primero para evitar conflictos
tmux kill-session -t $SESSION 2>/dev/null
 
# Crea la sesion con el primer cliente (no se attachea todavia, -d = detached)
tmux new-session -d -s $SESSION -n main "go run cmd/main.go client1 2 1"
 
# Divide la ventana en mas paneles y lanza cada cliente
tmux split-window -t $SESSION -h "go run cmd/main.go client2 2 2"
tmux split-window -t $SESSION -v "go run cmd/main.go client3 3 3"
 
tmux select-pane -t $SESSION:0.0
tmux split-window -t $SESSION -v "go run cmd/main.go client4 1 4"
 
tmux select-pane -t $SESSION:0.1
tmux split-window -t $SESSION -v "go run cmd/main.go client5 2 5"
 
# Acomoda los paneles en un grid parejo
tmux select-layout -t $SESSION tiled
 
# Te conecta a la sesion para que veas los 5 en vivo
tmux attach -t $SESSION
 