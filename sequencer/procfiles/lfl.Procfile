leader1: $SEQ_BIN -IP :8000 -leader -root -color 1
follower: sleep 1s && $SEQ_BIN -IP :8001 -parIP :8000 -color 1
leader2: sleep 2s && $SEQ_BIN -IP :8002 -parIP :8001 -leader -color 2