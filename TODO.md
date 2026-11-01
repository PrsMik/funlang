1. Убрать структрную типизацию в types/util.go: Equals/EqualsSoft для InterfaceType и ImplType.
2. В функции T стирается в MetaType (resolveType, checkFunctionLiteral, checkComptimeExpression).
3. ~~Утечка имён параметров функции во внешнее окружение.~~
4. ~~Неправильные восстановления curExpectedType в чекере.~~
5. Неправильная логика проверки ошибок во время comptime в checkComptimeExpression.
6. Типы массивов и мап с параметрами type берутся из чекера.
7. Проблемы с окружением в evalTypeField, evalTypeAccessExpression и в операторах | и &, if/switch.
8. evalInterfaceLiteral исполняет let в окружении вызывающего.
9. Ошибка в convertToTypeContext с восстановлением таблицы функций.
10. Паника на keyIndent.String() в checkFieldsInstantation и в evalFieldsInstantation.
11. ~~Баг в checkCallExpression при вызове функции без параметров с аргументами.~~
12. В IsAssignable: объединения сравниваются позиционно, нет ветки для ожидаемого & с не-& значением, неверно сравниваются параметры функций.
13.  Сужение res до типа из case не проверяет может ли тип быть вариантом объединения.
14. checkImplLiteral пишет в Fields фактический тип, а не объявленный
15. isImplementingType для | пишет ошибки левой ветки, даже если подойдёт правая.
16. TailCall возвращается из if/switch.
17. push требует строгого совпадения типов.