import sys, os
from ctypes import *

default_path = 'bin/'
default_fn   = 'pydc_client.so' 

def main(args):
    lib_path = os.path.join(args[1] if len(args) > 1 else default_path, default_fn) 
    pydc = cdll.LoadLibrary(lib_path)

    print(pydc.Initstore())
    print(pydc.Newcoll("testcoll".encode('utf-8')))
    print(pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcol".encode('utf-8'), 
        "int".encode('utf-8'),
    ))
    print(pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcolstr".encode('utf-8'), 
        "string".encode('utf-8'),
    ))
    print(pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcolflo".encode('utf-8'), 
        "float".encode('utf-8'),
    ))
    print(pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcolbool".encode('utf-8'), 
        "bool".encode('utf-8'),
    ))
    

    print('done')


if __name__=="__main__":
    main(sys.argv) 
