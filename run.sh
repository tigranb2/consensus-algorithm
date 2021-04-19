if [[ -z $5 ]]; then
  echo "Please specify the algorithm, number of nodes, broadcast delay/timeout, loss rate, and description"
else
  chmod 700 $1
    for i in {1..5}; do
    python3 start.py $1 $2 $3 $4
    python3 analysis.py $5
  done
  cat $5-data.txt
fi
