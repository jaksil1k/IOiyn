# IOiyn
before starting db run this sql statement

createSchemaQuery := "create database if not exists game CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"

if you are not Zaur or Alikhan, then it makes no sense for you to read text below.

будет Main page где есть последние игры, новые игры, самые популярные или "рейтинговые" и бесплатные игры как в магазин в стиме.

еще будет каталог где можно сортировать по году выпуска, рейтингу, цене и жанру.

еще будет кнопка как snippetCreate в книжке только где создаешь оффер для игры.

пока есть только эти методы в handlers: home, gameCreate, gameView, catalogView.

не знаю будет ли Admin panel и Autorization как успею пока можешь его не делать. Сам пока не решил, но приблизтельно
там будет две странички типо созданные игры и купленные или сделать одной библиотекой и там как то разделить, затем найстройки аккаунта изменить email, имя, пароль, 
nickname. не делай этого но возможно потом после файнала добавлю потверждение email, и восстановление пароля через email,
и рассылку(observer) можно сделать по приколу.

Пока сделал только до 4 чаптера по книжке. так что нету Dynamic HTML TEMPLATES, и нету базы данных(подумал что ты запаришься) но по желаению можем их довавить вместе 
в дискорде.
