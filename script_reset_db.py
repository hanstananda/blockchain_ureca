import subprocess
from timeit import default_timer as timer

start = timer()

addresses = [
'1N9skadjj8GZkUJkzCfEKwDCGJrt6cydsk',
'15DAKLs5bkUmTWAhTwnpFBqnMfpdakBXoC',
'141qmuD6Wh93Dg5vdqdDjGwTShuuiYBTZS',
'19XaVuXCxykqivr6KjhqVy6pgzZH2YivNT',
'1HiaSjZZoMP3s18edwWWZtwcs8QX2aZ9MK',
'14b7CHfHyQi3xaiSnxKSw7QbY7byjLj58e',
'1PcmXrssw54smjrfhuS1RufgNooYyBmEbv',
'15HyA5Gg2xRpDQPqfoNwYQ1w9BBChggktn',
'1QJ5P7B8F2PcYMYhaYtAXTJyzYzKdmzZTg',
'13WiDjGU1G5ASQNFKyoBfpaBka2qgvanN6',
'148Jd4NhdAxFz3GzKa4uSVeEXWjfoMrqxo',
'1GkM5TRDe59j4N3qrM4FzEYoazbFPYMTPV',
'121zhn7VbS9wcrikK5SN2JLhy4wUg6Luf9',
]

commands_notary = '''
del blockchain*.db
del blockchain*.db.lock
set NODE_ID=3000
echo %NODE_ID%
blockchain_ureca.exe createblockchain -address 121zhn7VbS9wcrikK5SN2JLhy4wUg6Luf9
'''

commands_generate = 'blockchain_ureca.exe generate -amount 200000 -offline -to '

commands_copy = '''
copy blockchain_3000.db blockchain_3001.db
copy blockchain_3000.db blockchain_3002.db
copy blockchain_3000.db blockchain_3003.db
'''

for i in range(len(addresses)):
    commands_notary += commands_generate + addresses[i] + '\n'

commands_notary += commands_copy
result = []

process_notary = subprocess.Popen("cmd", shell=True ,stdin=subprocess.PIPE, stdout=subprocess.PIPE,   stderr=subprocess.PIPE)
out, err = process_notary.communicate(commands_notary.encode('utf-8'))
print (out.decode(encoding='windows-1252'))
end = timer()
print(end-start)
