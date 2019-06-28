import subprocess
import time
from timeit import default_timer as timer

start = timer()

commands_node1 = '''
set NODE_ID=3001
'''

from_addresses = [
'1N9skadjj8GZkUJkzCfEKwDCGJrt6cydsk',
'15DAKLs5bkUmTWAhTwnpFBqnMfpdakBXoC',
'141qmuD6Wh93Dg5vdqdDjGwTShuuiYBTZS',
'19XaVuXCxykqivr6KjhqVy6pgzZH2YivNT',
'1HiaSjZZoMP3s18edwWWZtwcs8QX2aZ9MK',
'14b7CHfHyQi3xaiSnxKSw7QbY7byjLj58e',
]

addresses = [
    '1PcmXrssw54smjrfhuS1RufgNooYyBmEbv',
    '15HyA5Gg2xRpDQPqfoNwYQ1w9BBChggktn',
    '1QJ5P7B8F2PcYMYhaYtAXTJyzYzKdmzZTg',
    '13WiDjGU1G5ASQNFKyoBfpaBka2qgvanN6',
    '148Jd4NhdAxFz3GzKa4uSVeEXWjfoMrqxo',
    '1GkM5TRDe59j4N3qrM4FzEYoazbFPYMTPV',
    '121zhn7VbS9wcrikK5SN2JLhy4wUg6Luf9',
]

send_repeat = ["blockchain_ureca.exe send -from ",
               " -amount 1 -to "]


def copy_db():
    commands = "set NODE_ID=3002\n"
    if i % 100 == 0 and i > 0:
        print("i: ", i)
    commands += "copy blockchain_3000.db blockchain_3002.db" + '\n'
    process_node = subprocess.Popen("cmd", shell=True ,stdin=subprocess.PIPE, stdout=subprocess.PIPE,   stderr=subprocess.PIPE)
    process_node.communicate(commands.encode('utf-8'))


for t in range(1):
    print("t: ", t)
    for i in range(len(from_addresses)):
        commands_node1 = "set NODE_ID=3002\n"
        # if i % 1 == 0 and i > 0:
        print("i: ", i)
        commands_node1 += send_repeat[0] + from_addresses[i] + send_repeat[1] + addresses[9] + '\n'
        print("commands_node1: ", commands_node1)
        process_node1 = subprocess.Popen("cmd", shell=True ,stdin=subprocess.PIPE, stdout=subprocess.PIPE,   stderr=subprocess.PIPE)
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
