import sys, os
import marshal
from ctypes import *

default_path = 'bin/'
default_fn   = 'pydc_client.so' 

def main(args):
    lib_path = os.path.join(args[1] if len(args) > 1 else default_path, default_fn) 
    pydc = cdll.LoadLibrary(lib_path)

    test_tuple = (63728, 'stuff', 2.5, False)
    print(test_tuple)
    ttmd = marshal.dumps(test_tuple)

    print("entering go")
    pydc.TestFunc(
        ttmd,
        len(ttmd),
    )


if __name__=="__main__":
    main(sys.argv) 
