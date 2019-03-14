import random, time
from multiprocessing import Process, Pipe,cpu_count
import csv
from tqdm import tqdm

# Python implementation of Sequential and Parallel Quicksort
def main():
    maxRandom, totalNumbers = 1000, 1000
    Measurements = [ ['maxRandom', 'totalNumbers', 'SequentialTime', 'ParallelTime'] ]
    i = 1
    while totalNumbers < 15000000:
        create_list = [random.randint(1,maxRandom) for x in range(totalNumbers)]
        sequentialsortlist=create_list
        start = time.time()
        
        sorted2 = quicksort(sequentialsortlist)
        sequentialelapsed = time.time() - start

        n = cpu_count()
        pconn, cconn = Pipe()
        p = Process(target=quicksortParallel,
                    args=(create_list, cconn, n,))
        
        start = time.time()
        p.start()
        lst = pconn.recv()
        p.join()

        parallelelapsed = time.time() - start
        
        maxRandom+=150000
        totalNumbers+=150000
        i+=1
        Measurements.append([maxRandom,totalNumbers,sequentialelapsed, parallelelapsed])
        if i%10 == 0:
            print('Iteration number : ', i)

    with open("results_python.csv", "w") as f:
        writer = csv.writer(f)
        writer.writerows(Measurements)


def quicksort(lst):

    lesslist = []
    pivotList = []
    aheadlist = []
    if len(lst) <= 1:
        return lst
    else:
        pivot = lst[0]
        for i in lst:
            if i < pivot:
                lesslist.append(i)
            elif i > pivot:
                aheadlist.append(i)
            else:
                pivotList.append(i)
        lesslist = quicksort(lesslist)
        aheadlist = quicksort(aheadlist)
        return lesslist + pivotList + aheadlist

def quicksortParallel(lst, conn, procNum):
   lesslist = []
   pivotList = []
   aheadlist = []

   if procNum <= 0 or len(lst) <= 1:
       conn.send(quicksort(lst))
       conn.close()
       return
   else:
       pivot = lst[0]
       for i in lst:
           if i < pivot:
               lesslist.append(i)
           elif i > pivot:
               aheadlist.append(i)
           else:
               pivotList.append(i)

   pconnLeft, cconnLeft = Pipe()
   leftProc = Process(target=quicksortParallel,
                      args=(lesslist, cconnLeft, procNum - 1))
   pconnRight, cconnRight = Pipe()
   rightProc = Process(target=quicksortParallel,
                      args=(aheadlist, cconnRight, procNum - 1))
   leftProc.start()
   rightProc.start()
   conn.send(pconnLeft.recv()+pivotList + pconnRight.recv())
   conn.close()
   leftProc.join()
   rightProc.join()


if __name__ == '__main__':
    main()