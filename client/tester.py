from ctypes import *

def main():
    pystore = cdll.LoadLibrary('pydc_client.so')
    pystore.AddCol.argtypes = [c_char_p, c_char_p, c_char_p]

    print(pystore.Initstore())
    print(pystore.Newcoll("testcoll".encode('utf-8')))
    print(pystore.AddCol(
        "testcoll".encode('utf-8'),
        "testcol".encode('utf-8'), 
        "int".encode('utf-8'),
    ))
    print('done')


if __name__=="__main__":
    main() 
