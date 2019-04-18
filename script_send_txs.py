import subprocess
from timeit import default_timer as timer

start = timer()

commands_node1 = '''
export NODE_ID=3001
'''

send_repeat = "./blockchain_ureca send -from 13L7UYXjUCGUUKF5o4oExDFQnV6p3AkDoB " \
              "-to 1B84VWxLDwk2BBLnEhQioV1ZNJxxFmHpdA -amount 1\n"

for i in range(1):
    if i % 100 == 0:
        print(i)
    commands_node1 += send_repeat

# commands_node1 += "./blockchain_ureca startnode -port 9090\n"

# print(commands_node1)

process_node1 = subprocess.Popen('/bin/bash', stdin=subprocess.PIPE, stdout=subprocess.PIPE)
out, err = process_node1.communicate(commands_node1.encode('utf-8'))

# print(out)

end = timer()
print(end-start)
