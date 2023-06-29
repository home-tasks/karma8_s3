### Example of work
Execute
```
go run main.go
```

1. It will split `./example_inputs/in.txt` into 6 parts and store them asynchronously over 9 abstract storages (file providers, in reality) in `./out`.

2. A new server is added at runtime.

3. Then it will split `./example_inputs/in2.txt` into 6 parts and store them asynchronously over 10 "servers" (in `./out`).

4. Finally, it downloads asynchronously both files from those abstract storages and concatinates them in the main thread and prints to Stdout (for simplicity).

-------
The output looks similar to this:
```
servers: 9
8 bytes of 'example_inputs/in.txt' at storageId: 7 | path: out/7/0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt
8 bytes of 'example_inputs/in.txt' at storageId: 5 | path: out/5/0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt
8 bytes of 'example_inputs/in.txt' at storageId: 6 | path: out/6/0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt
8 bytes of 'example_inputs/in.txt' at storageId: 4 | path: out/4/0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt
8 bytes of 'example_inputs/in.txt' at storageId: 2 | path: out/2/0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt
8 bytes of 'example_inputs/in.txt' at storageId: 8 | path: out/8/0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt

A new server is added with the id: 10 

16 bytes of 'example_inputs/in2.txt' at storageId: 9 | path: out/9/4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt
16 bytes of 'example_inputs/in2.txt' at storageId: 8 | path: out/8/4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt
16 bytes of 'example_inputs/in2.txt' at storageId: 10 | path: out/10/4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt
16 bytes of 'example_inputs/in2.txt' at storageId: 3 | path: out/3/4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt
16 bytes of 'example_inputs/in2.txt' at storageId: 5 | path: out/5/4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt
15 bytes of 'example_inputs/in2.txt' at storageId: 6 | path: out/6/4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt


===== Content of '0ecfe268-160b-491f-ab02-a2ba4d6d8207_in.txt' is: =====
line111
line222
line333
line444
line555
line666

=====

===== Content of '4c3d496f-8a96-4da4-ab00-e20b4586d8a6_in2.txt' is: =====
this file is different from the previous one, it has everything in one line and the line is big
=====

(to remove file parts, execute: 'rm -r ./out')
```

### Key notes
1. When splitting the files, it always tries to keep the number of bytes in each part multiple of 8.
For example, for a file of 100 bytes, the file parts will be:
```
16, 16, 16, 16, 16, 20
```
> See examples in `./src/utils/split_size_test.go`

2. The abstract storage is represented as an interface, the real implementation of which is just a file storage. But pretty much anything could be use that implements `Storage`. The current implementation stores the files in to `./out` directory, but it is easy to create the one, that sends them to some outer servers.

3. The same goes for strategy that chooses the servers. The current one just chooses them randomly, but there might be other ones, for example, that can get the free space of the servers and store to the ones that have more space.
