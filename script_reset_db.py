import subprocess
from timeit import default_timer as timer

start = timer()

addresses = [
    '13XfCX8bLpdu8YgnXPD4BDeBC5RyvqBfPh',
    '14L3zLQWPiXM6hZXdfmgjET8crM52VJpXX',
    '1C4tyo8poeG1uFioZjtgnLZKotEUZFJyVh',
    '18Nt9jiYVjm2TxCTHNSeYquriaauh5wfux',
    '16uqNuajndwknbHSQw1cfTvSgsXxa5Vxi8',
    '1AqNL5SPcuWqUT1SjTEQ3WGDLfy47HK74c',
    '17aju9bJh3G7xC9PAkQ1j5czizA31rN77S',
    '1Ci67qmp8KerJA3zZhsDC7AcXz8RCZwbt',
    '1MzLjrr737WtVpubSGxN6CUECBD2vnQqef',
    '165KxLW2bFms5wtKs2sNQXfD8TLQrehGCT',
    '14RJHhG374XyuTLfZ48qRxUdxRLWj3BcA7',
]

commands_notary = '''
ls | grep -P "^blockchain_[0-9]{4}.db" | xargs -d "\\n" rm
export NODE_ID=3000
echo $NODE_ID
./blockchain_ureca createblockchain -address 1DAP5SpEFRuqUacbXFzsAjUFG3FPeQzDim
'''

commands_generate = './blockchain_ureca generate -amount 200000 -offline -to '

commands_copy = '''
cp blockchain_3000.db blockchain_3001.db
cp blockchain_3000.db blockchain_3002.db
cp blockchain_3000.db blockchain_3003.db
'''

for i in range(len(addresses)):
    commands_notary += commands_generate + addresses[i] + '\n'

commands_notary += commands_copy

process_notary = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
out, err = process_notary.communicate(commands_notary.encode('utf-8'))

end = timer()
print(end-start)
