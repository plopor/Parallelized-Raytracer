#!/bin/bash
#
#SBATCH --job-name=rt_benchmark 
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/bwluo/Documents/rt/benchmark
#SBATCH --partition=debug
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=4:00:00

module load golang/1.19

go run ../main.go
go run ../main.go
go run ../main.go
go run ../main.go
go run ../main.go

go run ../main.go 2
go run ../main.go 2
go run ../main.go 2
go run ../main.go 2
go run ../main.go 2

go run ../main.go 4
go run ../main.go 4
go run ../main.go 4
go run ../main.go 4
go run ../main.go 4

go run ../main.go 6
go run ../main.go 6
go run ../main.go 6
go run ../main.go 6
go run ../main.go 6

go run ../main.go 8
go run ../main.go 8
go run ../main.go 8
go run ../main.go 8
go run ../main.go 8

go run ../main.go 12
go run ../main.go 12
go run ../main.go 12
go run ../main.go 12
go run ../main.go 12

go run ../main.go 2 true false
go run ../main.go 2 true false
go run ../main.go 2 true false
go run ../main.go 2 true false
go run ../main.go 2 true false

go run ../main.go 4 true false
go run ../main.go 4 true false
go run ../main.go 4 true false
go run ../main.go 4 true false
go run ../main.go 4 true false

go run ../main.go 6 true false
go run ../main.go 6 true false
go run ../main.go 6 true false
go run ../main.go 6 true false
go run ../main.go 6 true false

go run ../main.go 8 true false
go run ../main.go 8 true false
go run ../main.go 8 true false
go run ../main.go 8 true false
go run ../main.go 8 true false

go run ../main.go 12 true false
go run ../main.go 12 true false
go run ../main.go 12 true false
go run ../main.go 12 true false
go run ../main.go 12 true false

go run ../main.go 2 true true
go run ../main.go 2 true true
go run ../main.go 2 true true
go run ../main.go 2 true true
go run ../main.go 2 true true

go run ../main.go 4 true true
go run ../main.go 4 true true
go run ../main.go 4 true true
go run ../main.go 4 true true
go run ../main.go 4 true true

go run ../main.go 6 true true
go run ../main.go 6 true true
go run ../main.go 6 true true
go run ../main.go 6 true true
go run ../main.go 6 true true

go run ../main.go 8 true true
go run ../main.go 8 true true
go run ../main.go 8 true true
go run ../main.go 8 true true
go run ../main.go 8 true true

go run ../main.go 12 true true
go run ../main.go 12 true true
go run ../main.go 12 true true
go run ../main.go 12 true true
go run ../main.go 12 true true

go run ../main.go 2 false true
go run ../main.go 2 false true
go run ../main.go 2 false true
go run ../main.go 2 false true
go run ../main.go 2 false true

go run ../main.go 4 false true
go run ../main.go 4 false true
go run ../main.go 4 false true
go run ../main.go 4 false true
go run ../main.go 4 false true

go run ../main.go 6 false true
go run ../main.go 6 false true
go run ../main.go 6 false true
go run ../main.go 6 false true
go run ../main.go 6 false true

go run ../main.go 8 false true
go run ../main.go 8 false true
go run ../main.go 8 false true
go run ../main.go 8 false true
go run ../main.go 8 false true

go run ../main.go 12 false true
go run ../main.go 12 false true
go run ../main.go 12 false true
go run ../main.go 12 false true
go run ../main.go 12 false true

python graph-script.py $SLURM_JOB_ID