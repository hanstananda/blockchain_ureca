import subprocess
import time
from timeit import default_timer as timer

start = timer()

commands_node1 = '''
export NODE_ID=3001
'''

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
    '13L7UYXjUCGUUKF5o4oExDFQnV6p3AkDoB',
]

send_repeat = ["./blockchain_ureca send -from ",
               " -amount 1 -to "]


def copy_db():
    commands = "export NODE_ID=3002\n"
    if i % 100 == 0 and i > 0:
        print("i: ", i)
    commands += "cp blockchain_3000.db blockchain_3002.db" + '\n'
    process_node = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
    process_node.communicate(commands.encode('utf-8'))


for t in range(100):
    print("t: ", t)
    for i in range(10):
        commands_node1 = "export NODE_ID=3002\n"
        if i % 100 == 0 and i > 0:
            print("i: ", i)
        commands_node1 += send_repeat[0] + addresses[i] + send_repeat[1] + addresses[10] + '\n'
        process_node1 = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
        out, err = process_node1.communicate(commands_node1.encode('utf-8'))
        time.sleep(1)
    # Make sure the new txs has been put into database
    time.sleep(1)
    copy_db()
    time.sleep(0.5)




# commands_node1 += "./blockchain_ureca startnode -port 9090\n"

# print(commands_node1)

process_node1 = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
out, err = process_node1.communicate(commands_node1.encode('utf-8'))

# print(out)

end = timer()
print(end-start)
