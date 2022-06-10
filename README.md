# SCP_BHA_FFA_Bacher

Набор пакетов, реализующий алгоритмы BHA и FFA, и инструмент проведения измерений. 

## Программы *cmd*

На основе пакетов реализуют несколько простых программ. Исполняемые файлы находятся в папке *bin*.

### *algoapp*
Использование алгоритма на тестовом файле/файлах. Файлы должны иметь структуру, аналогичную 65-файлам http://people.brunel.ac.uk/~mastjjb/jeb/orlib/scpinfo.html
**Использование:** 


**-a, —algo** string Name of algo [ffa|bh|ffapso|ffarms] (default "bh")


**--alpha** float Param alpha for FFA algos. (default 0.1)


**--betta** float Param betta for FFA algos. (default 1)


**--collapse** string Choose collapse for BHA [rand|stand] (default "rand")


**-c, —conf** string Name of config file for expt (without extension) (default "../../configs/config.yaml")


**--dist** string Choose dist_type for FFA manhat (default "euclid")


**--gamma** float Param gamma for FFA algos. (default 0.01)


**--move** string Choose move_type for FFA stand (default "best")


**--res** string File to save.


**--scsv** Save in csv style.


**--srcexpt** string Instance files mask.


**--stab** Save in table style.


**--stagn** float Param stang_persent for BHA algo. (default 1)



**Пример (для linux)**
```console
user@bar:~$ ./algoapp --conf "configs/config"  -a ffa —res res.txt —stab —srcexpt ../../OR/scp42.txt
```

### *scpdbapp*
Используется для загрузки в БД PostgreSQL экземпляра(ов). БД имеет следующую [структуру](https://i.ibb.co/dM3C3Lz/erd-drawio.png)


**-c, —conct** string Get connection string ['user= password= dbname= sslmode=']


**—dir** Source is directory.


**-m, —mtrx** string Name of table for matrix (default "matrix")


**-e, —mtrx_elem** string Name of table for matrix elements (default "matrix_element")


**-s, —source** string Source path file/dir


**—temp** Source by template.


**Пример (для linux)**
```console
user@bar:~$ ./scpdbapp -c "user=dev_ password=0000 dbname=scpdb sslmode=disable" -s "../../OR" —dir
```

### *pgexptapp*
Терминальный интерфейс для запуска на экземпляре(ах) из базы данных. Требует аналогичный параметр подключения к БД.
```console
user@bar:~$ ./pgexptapp -"user=dev_ password=0000 dbname=scpdb sslmode=disable"
```

### *configapp*
**Использование:**


**-d, —disrete** string Discretization method.


**-e, —expt** int Num of experiences iterations.


**-f, —file** string Give name of config file (without extension)


**-i, —iters** int Num of iterations for algo.


**-o, —opt** string Choose option create


**-p, —popsize** int Population size.


**-t, —transfer** string Transfer function.


#### Дополнение
Для сборки необходим установленный golang.
Для компиляции исходного файла из *cmd* - `go build file.go`.

### Краткое описание пакетов

#### scpmath
Cодержит в себе набор математических функций, который применяется реализованными алгоритмами для решения SCP.
#### scpfunc
Данный пакет содержит структуру SolutionRepairer, основными полями
являются параметры, для работы с матрицей покрытия. С помощью метода `RepairSolution()` проводится исправление переданного решения, путем покрытия незатронутых строк. С помощью метода `RemoveRedudancy()` проводится удаление избыточных столбцов. Метод `CheckSolution()` соответственно проверяет корректность решения и возвращает ответ (логическое значение и строку-ответ).
#### scpalgo
В данном пакете определен интерфейс `ScpSolver`, ключевым методом
которого является `Solve()`. Данный метод принимает общие для реализованных алгоритмов параметры, такие как размер популяции, максимальное число итераций, срез стоимостей столбцов и т.д. Определены структуры-алгоритмов, которые реализуют метод `Solve()`.

#### scpexpt
Данный пакет содержит структуры и функции, упрощающие процесс проведения экспериментов. Структура `ExptParams` предназначена для хранения и передачи параметров для эксперимента. Структура `ScpExptMaker` предназначена для запуска и сохранения/вывода результатов экспериментов.
