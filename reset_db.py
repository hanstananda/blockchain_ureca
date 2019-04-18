import subprocess
from timeit import default_timer as timer

start = timer()

commands_notary = '''
ls | grep -P "^blockchain_[0-9]{4}.db" | xargs -d "\\n" rm
export NODE_ID=3000
echo $NODE_ID
./blockchain_ureca createblockchain -address 1DAP5SpEFRuqUacbXFzsAjUFG3FPeQzDim
./blockchain_ureca generate -to 13L7UYXjUCGUUKF5o4oExDFQnV6p3AkDoB -amount 10000
cp blockchain_3000.db blockchain_3001.db
cp blockchain_3000.db blockchain_3002.db
cp blockchain_3000.db blockchain_3003.db
'''

process_notary = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
out, err = process_notary.communicate(commands_notary.encode('utf-8'))

end = timer()
print(end-start)
