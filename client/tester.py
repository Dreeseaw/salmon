import sys, os
import marshal
from ctypes import *

default_path = 'bin/'
default_fn   = 'pydc_client.so' 

def insert(dc, coll, tuppy):
    marshed = marshal.dumps(tuppy)
    dc.Insert(
        coll.encode('utf-8'),
        marshed,
        len(marshed),
    )

def main(args):
    lib_path = os.path.join(args[1] if len(args) > 1 else default_path, default_fn) 
    pydc = cdll.LoadLibrary(lib_path)

    collection = "testcoll"

    pydc.Initstore()
    pydc.Newcoll("testcoll".encode('utf-8'))
    pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcol".encode('utf-8'), 
        "int".encode('utf-8'),
    )
    pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcolstr".encode('utf-8'), 
        "string".encode('utf-8'),
    )
    pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcolflo".encode('utf-8'), 
        "float".encode('utf-8'),
    )
    pydc.AddCol(
        "testcoll".encode('utf-8'),
        "testcolbool".encode('utf-8'), 
        "bool".encode('utf-8'),
    )
    print(collection, " coll created")

    for i in range(0,100):
        insert(pydc, "testcoll", (235+i, 'stuff', 3.3, True))
    for i in range(0,100):
        insert(pydc, "testcoll", (235+i, 'more', 3.3, False))

    print("post-insert")

    selectors = marshal.dumps([
        ('testcol', 238),
        ('testcolstr', 'more'),
    ])
    pydc.Select(
        "testcoll".encode('utf-8'),
        selectors,
        len(selectors),
    )

    print('done')


if __name__=="__main__":
    main(sys.argv) 
