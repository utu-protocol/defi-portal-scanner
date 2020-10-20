#!/usr/bin/env python3

# Copyright © 2020 Andrea Giacobino <no.andrea@gmail.com>

# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:

# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.

# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

import argparse
import asyncio
import logging
# graphql
from python_graphql_client import GraphqlClient


version = "1.0.0"


logging.basicConfig(
    encoding='utf-8',
    level=logging.INFO
)


def pair_to_str(pair: object) -> str:
    """
        Uniswap pair to string
        """
    _id = pair.get("id")
    logging.debug(f"token pair {_id}")
    _t0 = pair.get("token0", {}).get("name")
    _t1 = pair.get("token1", {}).get("name")
    return f"{_t0}/{_t1}"


class UTUCLient(object):
    """
    Client to send data to the utu trust engine
    """

    def __init__(self, url: str):
        self.url = True

    async def consume(self, events_queue: asyncio.Queue):
        while True:
            value = await events_queue.get()
            if value is None:
                break
            self.post_edge(value)

    def post_edge(self, edge: str):
        # TODO: make the actual request
        print(edge)


class Uniswap(object):
    def __init__(self, endpoint: str):
        self.endpoint = endpoint
        self.client = GraphqlClient(endpoint=endpoint)

    def process_tx(self, t):
        async def push(v):
            await self.queue.put(v)

        logging.debug(">>>>> ", t)
        logging.debug("transaction id: ", t.get("id"))
        for e in t.get("swaps", []):
            # from to
            sub = e.get("sender")
            obj = e.get("to")
            usd = float(e.get("amountUSD"))
            # pair
            p_name = pair_to_str(e.get("pair", {}))
            # sample record
            log = f"<TRADE> from:{sub} to:{obj} usd:{usd:10.2f} pair:{p_name}"
            push(log)

        for e in t.get("mint", []):
            to = e.get("to")
            fto = e.get("fee_to")
            usd = float(e.get("amountUSD"))
            p_name = pair_to_str(e.get("pair", {}))
            # sample record
            log = f"<TRADE> to:{to} fee:{fto} usd:{usd:10.2f} pair:{p_name}"
            push(log)

    async def produce(self, events_queue: asyncio.Queue):
        self.queue = events_queue
        query = '''transactions{
            id
            blockNumber
            timestamp
            swaps {
                id
                pair {
                    id
                    token0 {
                        name
                    }
                    token1 {
                        name
                    }
                }
                sender
                to
                amountUSD
            }
            mints {
                id
                liquidity
                feeTo
                to
                pair {
                    id
                    token0 {
                        name
                    }
                    token1 {
                        name
                    }
                }
            }
            burns {
                id
            }
        }'''
        query = "subscription { %s }" % query.replace("\n", " ")
        await self.client.subscribe(query=query, handle=print)


def graphql_query(endpoint: str, query: str) -> object:
    """
    execute a graphql query
    """
    try:

        query = "query { %s }" % query.replace("\n", " ")
        logging.debug(query)
        client = GraphqlClient(endpoint=endpoint)
        data = client.execute(query=query)
        logging.debug(data)
        return data.get("data", {})
    except Exception as err:
        print(f"Error querying: {err}")
        return {}


async def run_uniswap(endpoint: str, queue: asyncio.Queue):
    query = '''transactions(first: 100){
            id
            blockNumber
            timestamp
            swaps {
                id
                pair {
                    id
                    token0 {
                        name
                    }
                    token1 {
                        name
                    }
                }
                sender
                to
                amountUSD
            }
            mints {
                id
                liquidity
                feeTo
                to
                pair {
                    id
                    token0 {
                        name
                    }
                    token1 {
                        name
                    }
                }
            }
            burns {
                id
            }
        }'''

    # retrieve the transactions
    r = graphql_query(endpoint, query)
    txs = r.get("transactions", [])
    # go fetch the transactions
    for t in txs:
        logging.debug("transaction id: ", t["id"])
        for e in t.get("swaps", []):
            # from to
            sub = e.get("sender")
            obj = e.get("to")
            usd = float(e.get("amountUSD"))
            # pair
            p_name = pair_to_str(e.get("pair", {}))
            # sample record
            log = f"<TRADE> from:{sub} to:{obj} usd:{usd:10.2f} pair:{p_name}"
            await queue.put(log)

        for e in t.get("mint", []):
            to = e.get("to")
            fto = e.get("fee_to")
            usd = float(e.get("amountUSD"))
            p_name = pair_to_str(e.get("pair", {}))
            # sample record
            log = f"<TRADE> to:{to} fee:{fto} usd:{usd:10.2f} pair:{p_name}"
            await queue.put(log)

    await queue.put(None)


def cmd_uniswap(args):
    """
    Export data from graphql endpoint
    """

    # output folder
    # os.makedirs(args.output, exist_ok=True)

    # init the utu client
    utu = UTUCLient("https://ututrls")
    uni = Uniswap(args.graphql)

    # push the actions into the queue
    loop = asyncio.get_event_loop()
    actions = asyncio.Queue()
    from_eth = run_uniswap(args.graphql, actions)  # sync version
    # from_eth = uni.produce(actions)  # sync version
    to_utu = utu.consume(actions)
    loop.run_until_complete(asyncio.gather(from_eth, to_utu))
    loop.close()


def cmd_version(args):
    """
    Print the version and exit
    """
    print(f"download v{version}")


def main():
    commands = [
        {
            'name': 'uniswap',
            'help': 'Export a Thing messages into separate files to be analized with mallet',
            'target': cmd_uniswap,
            'opts': [
                {
                    "names": ["--debug", "-d"],
                    "help": "enable debug logging",
                    "default": False,
                    "action": "store_true",
                },
                {
                    "names": ["-o", "--output"],
                    "help": "the output folder",
                    "default": "rawdata"
                },
                {
                    "names": ["--graphql", "-g"],
                    "help": "the graphql endpoint",
                    "default": "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2"
                },

            ]
        },
        {
            'name': 'version',
            'help': 'Print the version and exit',
            'target': cmd_version,
            'opts': [
            ]
        },

    ]
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers()
    subparsers.required = True
    subparsers.dest = 'command'
    # register all the commands
    for c in commands:
        subparser = subparsers.add_parser(c['name'], help=c['help'])
        subparser.set_defaults(func=c['target'])
        # add the sub arguments
        for sa in c.get('opts', []):
            subparser.add_argument(*sa['names'],
                                   help=sa['help'],
                                   action=sa.get('action'),
                                   default=sa.get('default'))

    # parse the arguments
    args = parser.parse_args()
    if hasattr(args, "debug") and args.debug:
        logging.root.setLevel(logging.DEBUG)
    # call the function
    args.func(args)


if __name__ == "__main__":
    main()
