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

mid = timer()
print(mid-start)

# print(out)
commands_node1 = '''
export NODE_ID=3001
'''

send_repeat = "./blockchain_ureca send -from 13L7UYXjUCGUUKF5o4oExDFQnV6p3AkDoB " \
              "-to 1B84VWxLDwk2BBLnEhQioV1ZNJxxFmHpdA -amount 1\n"

for i in range(5000):
    commands_node1 += send_repeat

# commands_node1 += "./blockchain_ureca startnode -port 9090\n"

# print(commands_node1)

process_node1 = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
out, err = process_node1.communicate(commands_node1.encode('utf-8'))

# print(out)

end = timer()
print(end-start)
