package genres

const GENRES_XML = `

<!-- loaded from ftp://ftp.fictionbook.org/pub/genres/genres_transfer_utf8.zip and flibgolite -->
<?xml version="1.0" encoding="utf-8"?>
<fbgenrestransfer xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="http://alexs.ru/fb2/genrestable/GT.xsd">
    <genre value="sf">
        <root-descr lang="en" genre-title="SF, Fantasy" detailed="Alternate History, Cyberpunk, Fantasy, Science Fiction"/>
        <root-descr lang="ru" genre-title="Фантастика, Фэнтэзи" detailed="НФ, Фэнтези, Мистика, Киберпанк"/>
        <root-descr lang="uk" genre-title="Фантастика, Фентезі" detailed="НФ, Фентезі, Містика, Кіберпанк"/>
        <subgenres>
            <subgenre value="sf_history">
                <genre-descr lang="en" title="Alternative history"/>
                <genre-descr lang="ru" title="Альтернативная история"/>
                <genre-descr lang="uk" title="Альтернативна історія"/>
                <genre-alt value="fantasy_alt_hist" format="fb2.0"/>
                <genre-alt value="historical_fantasy" format="flibgolite"/>
                <genre-alt value="popadanec" format="flibgolite"/>
                <genre-alt value="popadantsy" format="flibgolite"/>
            </subgenre>
            <subgenre value="sf_action">
                <genre-descr lang="en" title="Action SF"/>
                <genre-descr lang="ru" title="Боевая Фантастика"/>
                <genre-descr lang="uk" title="Бойова Фантастика"/>
                <genre-alt value="fantasy_fight" format="flibgolite"/>
                <genre-alt value="fantasy_action" format="flibgolite"/>
            </subgenre>
            <subgenre value="sf_epic">
                <genre-descr lang="en" title="Epic SF"/>
                <genre-descr lang="ru" title="Эпическая Фантастика"/>
                <genre-descr lang="uk" title="Епічна Фантастика"/>
            </subgenre>
            <subgenre value="sf_heroic">
                <genre-descr lang="en" title="Heroic SF"/>
                <genre-descr lang="ru" title="Героическая Фантастика"/>
                <genre-descr lang="uk" title="Героїчна Фантастика"/>
            </subgenre>
            <subgenre value="sf_detective">
                <genre-descr lang="en" title="Detective SF"/>
                <genre-descr lang="ru" title="Детективная Фантастика"/>
                <genre-descr lang="uk" title="Детективна Фантастика"/>
            </subgenre>
            <subgenre value="sf_cyberpunk">
                <genre-descr lang="en" title="Cyberpunk"/>
                <genre-descr lang="ru" title="Киберпанк"/>
                <genre-descr lang="uk" title="Кіберпанк"/>
                <genre-alt value="sf_cyber_punk" format="fb2.0"/>
            </subgenre>
            <subgenre value="sf_space">
                <genre-descr lang="en" title="Space SF"/>
                <genre-descr lang="ru" title="Космическая Фантастика"/>
                <genre-descr lang="uk" title="Космічна Фантастика"/>
            </subgenre>
            <subgenre value="sf_social">
                <genre-descr lang="en" title="Social SF"/>
                <genre-descr lang="ru" title="Социальная Фантастика"/>
                <genre-descr lang="uk" title="Соціальна Фантастика"/>
            </subgenre>
            <subgenre value="sf_horror">
                <genre-descr lang="en" title="Horror & Mystic"/>
                <genre-descr lang="ru" title="Ужасы и Мистика"/>
                <genre-descr lang="uk" title="Жахи та Містика"/>
                <genre-alt value="gay_mystery" format="fb2.0"/>
                <genre-alt value="horror" format="fb2.0"/>
                <genre-alt value="horror_antology" format="fb2.0"/>
                <genre-alt value="horror_british" format="fb2.0"/>
                <genre-alt value="horror_fantasy" format="fb2.0"/>
                <genre-alt value="horror_erotic" format="fb2.0"/>
                <genre-alt value="horror_ghosts" format="fb2.0"/>
                <genre-alt value="horror_graphic" format="fb2.0"/>
                <genre-alt value="horror_occult" format="fb2.0"/>
                <genre-alt value="horror_ref" format="fb2.0"/>
                <genre-alt value="horror_usa" format="fb2.0"/>
                <genre-alt value="horror_vampires" format="fb2.0"/>
                <genre-alt value="teens_horror" format="fb2.0"/>
                <genre-alt value="sf_mystic" format="flibgolite"/>
            </subgenre>
            <subgenre value="sf_humor">
                <genre-descr lang="en" title="Humor SF"/>
                <genre-descr lang="ru" title="Юмористическая Фантастика"/>
                <genre-descr lang="uk" title="Гумористична Фантастика"/>
                <genre-alt value="humor_fantasy" format="flibgolite"/>
            </subgenre>
            <subgenre value="sf_fantasy">
                <genre-descr lang="en" title="Fantasy"/>
                <genre-descr lang="ru" title="Фэнтези"/>
                <genre-descr lang="uk" title="Фентезі"/>
                <genre-alt value="romance_fantasy" format="fb2.0"/>
                <genre-alt value="romance_sf" format="fb2.0"/>
                <genre-alt value="romance_time_travel" format="fb2.0"/>
                <genre-alt value="foreign_fantasy" format="flibgolite"/>
                <genre-alt value="russian_fantasy" format="flibgolite"/>
                <genre-alt value="city_fantasy" format="flibgolite"/>
                <genre-alt value="sf_fantasy_city" format="flibgolite"/>
                <genre-alt value="magician_book" format="flibgolite"/>
                <genre-alt value="dragon_fantasy" format="flibgolite"/>
                <genre-alt value="fantasy" format="flibgolite"/>
            </subgenre>
            <subgenre value="sf">
                <genre-descr lang="en" title="Science Fiction"/>
                <genre-descr lang="ru" title="Научная Фантастика"/>
                <genre-descr lang="uk" title="Наукова Фантастика"/>
                <genre-alt value="gaming" format="fb2.0"/>
                <genre-alt value="sf_writing" format="fb2.0"/>
                <genre-alt value="foreign_sf" format="flibgolite"/>
                <genre-alt value="sci_fi" format="flibgolite"/>
            </subgenre>
            <subgenre value="child_sf">
                <genre-descr lang="en" title="Science Fiction for Kids"/>
                <genre-descr lang="ru" title="Детская Фантастика"/>
                <genre-descr lang="uk" title="Дитяча Фантастика"/>
                <genre-alt value="teens_sf" format="fb2.0"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="love_sf"> 
                <genre-descr lang="en" title="Love fantasies, love fiction novels"/>
                <genre-descr lang="ru" title="Любовное фэнтези, любовно-фантастические романы"/>
                <genre-descr lang="uk" title="Любовне фентезі, любовно-фантастичні романи"/>
                <genre-alt value="love_fantasy" format="flibgolite"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="sf_litrpg">
                <genre-descr lang="en" title="Literary Role Playing Game"/>
                <genre-descr lang="ru" title="ЛитРПГ (литературная RPG)"/>
                <genre-descr lang="uk" title="ЛітРПГ (літературна RPG)"/>
                <genre-alt value="litrpg" format="flibgolite"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="sf_postapocalyptic">
                <genre-descr lang="en" title="Postapocalypse"/>
                <genre-descr lang="ru" title="Постапокалипсис"/>
                <genre-descr lang="uk" title="Постапокаліпсис"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="sf_stimpank">
                <genre-descr lang="en" title="Steampunk"/>
                <genre-descr lang="ru" title="Стимпанк"/>
                <genre-descr lang="uk" title="Стімпанк"/>
            </subgenre>
    </genre>
    <genre value="detective">
        <root-descr lang="en" genre-title="Detectives, Thrillers" detailed="Police Stories, Ironical, Espionage, Crime"/>
        <root-descr lang="ru" genre-title="Детективы, Боевики" detailed="Полицейские, иронические, шпионские, криминальные"/>
        <root-descr lang="uk" genre-title="Детективи, Бойовики" detailed="Поліцейські, іронічні, шпигунські, кримінальні"/>
        <subgenres>
            <subgenre value="det_classic">
                <genre-descr lang="en" title="Classical Detective"/>
                <genre-descr lang="ru" title="Классический Детектив"/>
                <genre-descr lang="uk" title="Класичний Детектив"/>
            </subgenre>
            <subgenre value="det_police">
                <genre-descr lang="en" title="Police Stories"/>
                <genre-descr lang="ru" title="Полицейский Детектив"/>
                <genre-descr lang="uk" title="Поліцейський Детектив"/>
                <genre-alt value="thriller_police" format="fb2.0"/>
            </subgenre>
            <subgenre value="det_action">
                <genre-descr lang="en" title="Action"/>
                <genre-descr lang="ru" title="Боевики"/>
                <genre-descr lang="uk" title="Бойовики"/>
            </subgenre>
            <subgenre value="det_irony">
                <genre-descr lang="en" title="Ironical Detective"/>
                <genre-descr lang="ru" title="Иронический Детектив"/>
                <genre-descr lang="uk" title="Іронічний детектив"/>
            </subgenre>
            <subgenre value="det_history">
                <genre-descr lang="en" title="Historical Detective"/>
                <genre-descr lang="ru" title="Исторический Детектив"/>
                <genre-descr lang="uk" title="Історичний детектив"/>
            </subgenre>
            <subgenre value="det_espionage">
                <genre-descr lang="en" title="Espionage Detective"/>
                <genre-descr lang="ru" title="Шпионский Детектив"/>
                <genre-descr lang="uk" title="Шпигунський Детектив"/>
            </subgenre>
            <subgenre value="det_crime">
                <genre-descr lang="en" title="Crime Detective"/>
                <genre-descr lang="ru" title="Криминальный Детектив"/>
                <genre-descr lang="uk" title="Кримінальний детектив"/>
            </subgenre>
            <subgenre value="det_political">
                <genre-descr lang="en" title="Political Detective"/>
                <genre-descr lang="ru" title="Политический Детектив"/>
                <genre-descr lang="uk" title="Політичний детектив"/>
            </subgenre>
            <subgenre value="det_maniac">
                <genre-descr lang="en" title="Maniacs"/>
                <genre-descr lang="ru" title="Маньяки"/>
                <genre-descr lang="uk" title="Маньякі"/>
            </subgenre>
            <subgenre value="det_hard">
                <genre-descr lang="en" title="Hard-boiled Detective"/>
                <genre-descr lang="ru" title="Крутой Детектив"/>
                <genre-descr lang="uk" title="Крутий Детектив"/>
            </subgenre>
            <subgenre value="thriller">
                <genre-descr lang="en" title="Thrillers"/>
                <genre-descr lang="ru" title="Триллеры"/>
                <genre-descr lang="uk" title="Трилери"/>
                <genre-alt value="thriller_mystery" format="fb2.0"/>
            </subgenre>
            <subgenre value="detective">
                <genre-descr lang="en" title="Detective"/>
                <genre-descr lang="ru" title="Детектив"/>
                <genre-descr lang="uk" title="Детектив"/>
                <genre-alt value="mystery" format="fb2.0"/>
                <genre-alt value="foreign_detective" format="flibgolite"/>
            </subgenre>
            <subgenre value="sf_detective">
                <genre-descr lang="en" title="Detective SF"/>
                <genre-descr lang="ru" title="Детективная Фантастика"/>
                <genre-descr lang="uk" title="Детективна Фантастика"/>
                <genre-alt value="teens_mysteries" format="fb2.0"/>
            </subgenre>
            <subgenre value="child_det">
                <genre-descr lang="en" title="Children's Action"/>
                <genre-descr lang="ru" title="Детские Остросюжетные"/>
                <genre-descr lang="uk" title="Дитячі гостросюжетні"/>
            </subgenre>
            <subgenre value="love_detective">
                <genre-descr lang="en" title="Detective Romance"/>
                <genre-descr lang="ru" title="Остросюжетные Любовные Романы"/>
                <genre-descr lang="uk" title="Гостросюжетні Любовні Романи"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="prose">
        <root-descr lang="en" genre-title="Prose" detailed="Classical, History, Contemporary"/>
        <root-descr lang="ru" genre-title="Проза" detailed="Классика, историческая, современная"/>
        <root-descr lang="uk" genre-title="Проза" detailed="Класика, історична, сучасна"/>
        <subgenres>
            <subgenre value="prose_classic">
                <genre-descr lang="en" title="Classics Prose"/>
                <genre-descr lang="ru" title="Классическая Проза"/>
                <genre-descr lang="uk" title="Класична Проза"/>
                <genre-alt value="literature" format="fb2.0"/>
                <genre-alt value="literature_books" format="fb2.0"/>
                <genre-alt value="literature_british" format="fb2.0"/>
                <genre-alt value="literature_classics" format="fb2.0"/>
                <genre-alt value="literature_drama" format="fb2.0"/>
                <genre-alt value="literature_essay" format="fb2.0"/>
                <genre-alt value="literature_antology" format="fb2.0"/>
                <genre-alt value="literature_saga" format="fb2.0"/>
                <genre-alt value="literature_short" format="fb2.0"/>
                <genre-alt value="literature_usa" format="fb2.0"/>
                <genre-alt value="literature_world" format="fb2.0"/>
                <genre-alt value="prose" format="flibgolite"/>
                <genre-alt value="short_story" format="flibgolite"/>
                <genre-alt value="literature_20" format="flibgolite"/>
                <genre-alt value="literature_19" format="flibgolite"/>
                <genre-alt value="literature_18" format="flibgolite"/>
                <genre-alt value="foreign_prose" format="flibgolite"/>
            </subgenre>
            <subgenre value="prose_history">
                <genre-descr lang="en" title="Historical Prose"/>
                <genre-descr lang="ru" title="Историческая Проза"/>
                <genre-descr lang="uk" title="Історична Проза"/>
                <genre-alt value="literature_history" format="fb2.0"/>
                <genre-alt value="literature_critic" format="fb2.0"/>
            </subgenre>
            <subgenre value="prose_contemporary">
                <genre-descr lang="en" title="Contemporary Prose"/>
                <genre-descr lang="ru" title="Современная Проза"/>
                <genre-descr lang="uk" title="Сучасна Проза"/>
                <genre-alt value="literature_political" format="fb2.0"/>
                <genre-alt value="literature_war" format="fb2.0"/>
                <genre-alt value="ref_writing" format="fb2.0"/>
                <genre-alt value="foreign_contemporary" format="flibgolite"/>
                <genre-alt value="russian_contemporary" format="flibgolite"/>
            </subgenre>
            <subgenre value="prose_counter">
                <genre-descr lang="en" title="Counterculture"/>
                <genre-descr lang="ru" title="Контркультура"/>
                <genre-descr lang="uk" title="Контркультура"/>
                <genre-alt value="literature_gay" format="fb2.0"/>
            </subgenre>
            <subgenre value="prose_rus_classsic">
                <genre-descr lang="en" title="Russian Classics"/>
                <genre-descr lang="ru" title="Русская Классика"/>
                <genre-descr lang="uk" title="Російська Класика"/>
                <genre-alt value="literature_rus_classsic" format="fb2.0"/>
                <genre-alt value="prose_rus_classic" format="fb2.0"/>
            </subgenre>
            <subgenre value="prose_su_classics">
                <genre-descr lang="en" title="Soviet Classics"/>
                <genre-descr lang="uk" title="Радянська Класика"/>
                <genre-alt value="literature_su_classics" format="fb2.0"/>
            </subgenre>
            <subgenre value="humor_prose">
                <genre-descr lang="en" title="Humor Prose"/>
                <genre-descr lang="ru" title="Юмористическая Проза"/>
                <genre-descr lang="uk" title="Гумористична Проза"/>
            </subgenre>
            <subgenre value="child_prose">
                <genre-descr lang="en" title="Children's Prose"/>
                <genre-descr lang="ru" title="Детская Проза"/>
                <genre-descr lang="uk" title="Дитяча Проза"/>
                <genre-alt value="teens_literature" format="fb2.0"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="network_literature">
                <genre-descr lang="en" title="Selfpublished, online literature"/>
                <genre-descr lang="ru" title="Самиздат, сетевая литература"/>
                <genre-descr lang="uk" title="Самвидав, мережева література"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="aphorisms">
                <genre-descr lang="en" title="Aphorisms, quotes"/>
                <genre-descr lang="ru" title="Афоризмы, цитаты"/>
                <genre-descr lang="uk" title="Афоризми, цитати"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="prose_military">
                <genre-descr lang="en" title="War prose"/>
                <genre-descr lang="ru" title="Проза о войне"/>
                <genre-descr lang="uk" title="Проза про війну"/>
            </subgenre>
            <!-- flibgolite -->
            <subgenre value="fanfiction">
                <genre-descr lang="en" title="Fan Fiction"/>
                <genre-descr lang="ru" title="Фанфик"/>
                <genre-descr lang="uk" title="Фанфік"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="love">
        <root-descr lang="en" genre-title="Romance" detailed="Historical, Contemporary, Detective"/>
        <root-descr lang="ru" genre-title="Любовные романы" detailed="Исторические, современные, остросюжетные"/>
        <root-descr lang="uk" genre-title="Любовні романи" detailed="Історичні, сучасні, гостросюжетні"/>
        <subgenres>
            <subgenre value="love_contemporary">
                <genre-descr lang="en" title="Contemporary Romance"/>
                <genre-descr lang="ru" title="Современные Любовные Романы"/>
                <genre-descr lang="uk" title="Сучасні Любовні Романи"/>
                <genre-alt value="romance" format="fb2.0"/>
                <genre-alt value="romance_multicultural" format="fb2.0"/>
                <genre-alt value="romance_series" format="fb2.0"/>
                <genre-alt value="romance_anthologies" format="fb2.0"/>
                <genre-alt value="romance_contemporary" format="fb2.0"/>
                <genre-alt value="literature_women" format="fb2.0"/>
                <genre-alt value="romance_romantic_suspense" format="fb2.0"/>
            </subgenre>
            <subgenre value="love_history">
                <genre-descr lang="en" title="Historical Romance"/>
                <genre-descr lang="ru" title="Исторические Любовные Романы"/>
                <genre-descr lang="uk" title="Історичні Любовні Романи"/>
                <genre-alt value="romance_regency" format="fb2.0"/>
                <genre-alt value="romance_historical" format="fb2.0"/>
            </subgenre>
            <subgenre value="love_detective">
                <genre-descr lang="en" title="Detective Romance"/>
                <genre-descr lang="ru" title="Остросюжетные Любовные Романы"/>
                <genre-descr lang="uk" title="Гостросюжетні Любовні Романи"/>
            </subgenre>
            <subgenre value="love_short">
                <genre-descr lang="en" title="Short Romance"/>
                <genre-descr lang="ru" title="Короткие Любовные Романы"/>
                <genre-descr lang="uk" title="Короткі Любовні Романи"/>
            </subgenre>
            <subgenre value="love_erotica">
                <genre-descr lang="en" title="Erotica"/>
                <genre-descr lang="ru" title="Эротика"/>
                <genre-descr lang="uk" title="Еротика"/>
                <genre-alt value="literature_erotica" format="fb2.0"/>
                <genre-alt value="love_hard" format="flibgolite"/>
            </subgenre>
            <subgenre value="love">
                <genre-descr lang="en" title="Romance"/>
                <genre-descr lang="ru" title="Любовные Романы"/>
                <genre-descr lang="uk" title="Любовні Романи"/>
                <genre-alt value="foreign_love" format="flibgolite"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="adventure">
        <root-descr lang="en" genre-title="Adventure" detailed="Western, Historical, Sea"/>
        <root-descr lang="ru" genre-title="Приключения" detailed="Вестерны, исторические, морские"/>
        <root-descr lang="uk" genre-title="Пригоди" detailed="Вестерни, історичні, морські"/>
        <subgenres>
            <subgenre value="adv_western">
                <genre-descr lang="en" title="Western"/>
                <genre-descr lang="ru" title="Вестерны"/>
                <genre-descr lang="uk" title="Вестерни"/>
                <genre-alt value="literature_western" format="fb2.0"/>
            </subgenre>
            <subgenre value="adv_history">
                <genre-descr lang="en" title="History"/>
                <genre-descr lang="ru" title="Исторические Приключения"/>
                <genre-descr lang="uk" title="Історичні пригоди"/>
            </subgenre>
            <subgenre value="adv_indian">
                <genre-descr lang="en" title="Indians"/>
                <genre-descr lang="ru" title="Приключения: Индейцы"/>
                <genre-descr lang="uk" title="Пригоди: Індіанці"/>
            </subgenre>
            <subgenre value="adv_maritime">
                <genre-descr lang="en" title="Maritime Fiction"/>
                <genre-descr lang="ru" title="Морские Приключения"/>
                <genre-descr lang="uk" title="Морські пригоди"/>
                <genre-alt value="literature_sea" format="fb2.0"/>
            </subgenre>
            <subgenre value="adv_geo">
                <genre-descr lang="en" title="Travel & Geography"/>
                <genre-descr lang="ru" title="Путешествия и География"/>
                <genre-descr lang="uk" title="Подорожі та Географія"/>
                <genre-alt value="gay_travel" format="fb2.0"/>
                <genre-alt value="outdoors_travel" format="fb2.0"/>
                <genre-alt value="travel" format="fb2.0"/>
                <genre-alt value="travel_africa" format="fb2.0"/>
                <genre-alt value="travel_asia" format="fb2.0"/>
                <genre-alt value="travel_australia" format="fb2.0"/>
                <genre-alt value="travel_canada" format="fb2.0"/>
                <genre-alt value="travel_caribbean" format="fb2.0"/>
                <genre-alt value="travel_europe" format="fb2.0"/>
                <genre-alt value="travel_guidebook_series" format="fb2.0"/>
                <genre-alt value="travel_lat_am" format="fb2.0"/>
                <genre-alt value="travel_middle_east" format="fb2.0"/>
                <genre-alt value="travel_polar" format="fb2.0"/>
                <genre-alt value="travel_spec" format="fb2.0"/>
                <genre-alt value="travel_usa" format="fb2.0"/>
                <genre-alt value="travel_rus" format="fb2.0"/>
                <genre-alt value="travel_ex_ussr" format="fb2.0"/>
            </subgenre>
            <subgenre value="adv_animal">
                <genre-descr lang="en" title="Nature & Animals"/>
                <genre-descr lang="ru" title="Природа и Животные"/>
                <genre-descr lang="uk" title="Природа та Тварини"/>
                <genre-alt value="child_animals" format="fb2.0"/>
            </subgenre>
            <subgenre value="adventure">
                <genre-descr lang="en" title="Misk Adventures"/>
                <genre-descr lang="ru" title="Приключения: Прочее"/>
                <genre-descr lang="uk" title="Пригоди: Інше"/>
                <genre-alt value="literature_adv" format="fb2.0"/>
                <genre-alt value="literature_men_advent" format="fb2.0"/>
            </subgenre>
            <subgenre value="child_adv">
                <genre-descr lang="en" title="Adventures for Kids"/>
                <genre-descr lang="ru" title="Детские Приключения"/>
                <genre-descr lang="uk" title="Дитячі Пригоди"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="children">
        <root-descr lang="en" genre-title="Children's" detailed="Fairy Tales, Fantasy, Detectives..."/>
        <root-descr lang="ru" genre-title="Детское" detailed="Сказки, фантастика, детективы..."/>
        <root-descr lang="uk" genre-title="Дитяче" detailed="Казки, фантастика, детективи..."/>
        <subgenres>
            <subgenre value="child_tale">
                <genre-descr lang="en" title="Fairy Tales"/>
                <genre-descr lang="ru" title="Сказки"/>
                <genre-descr lang="uk" title="Казки"/>
                <genre-alt value="child_3" format="fb2.0"/>
                <genre-alt value="literature_fairy" format="fb2.0"/>
            </subgenre>
            <subgenre value="child_verse">
                <genre-descr lang="en" title="Verses"/>
                <genre-descr lang="ru" title="Детские Стихи"/>
                <genre-descr lang="uk" title="Дитячі Вірші"/>
            </subgenre>
            <subgenre value="child_prose">
                <genre-descr lang="en" title="Prose for Kids"/>
                <genre-descr lang="ru" title="Детская Проза"/>
                <genre-descr lang="uk" title="Дитяча Проза"/>
                <genre-alt value="child_4" format="fb2.0"/>
                <genre-alt value="child_9" format="fb2.0"/>
                <genre-alt value="child_history" format="fb2.0"/>
                <genre-alt value="child_characters" format="fb2.0"/>
            </subgenre>
            <subgenre value="child_sf">
                <genre-descr lang="en" title="Science Fiction for Kids"/>
                <genre-descr lang="ru" title="Детская Фантастика"/>
                <genre-descr lang="uk" title="Дитяча Фантастика"/>
            </subgenre>
            <subgenre value="child_det">
                <genre-descr lang="en" title="Detectives & Thrillers"/>
                <genre-descr lang="ru" title="Детские Остросюжетные"/>
                <genre-descr lang="uk" title="Дитячі Гостросюжетні"/>
            </subgenre>
            <subgenre value="child_adv">
                <genre-descr lang="en" title="Adventures for Kids"/>
                <genre-descr lang="ru" title="Детские Приключения"/>
                <genre-descr lang="uk" title="Дитячі Пригоди"/>
                <genre-alt value="teens_history" format="fb2.0"/>
                <genre-alt value="teens_series" format="fb2.0"/>
            </subgenre>
            <subgenre value="child_education">
                <genre-descr lang="en" title="Education for Kids"/>
                <genre-descr lang="ru" title="Детская образовательная литература"/>
                <genre-descr lang="uk" title="Дитяча освітня література"/>
                <genre-alt value="child_edu" format="fb2.0"/>
                <genre-alt value="child_nature" format="fb2.0"/>
            </subgenre>
            <subgenre value="children">
                <genre-descr lang="en" title="For Kids: Misk"/>
                <genre-descr lang="ru" title="Детское: Прочее"/>
                <genre-descr lang="uk" title="Дитяче: Інше"/>
                <genre-alt value="child_art" format="fb2.0"/>
                <genre-alt value="child_obsessions" format="fb2.0"/>
                <genre-alt value="child_people" format="fb2.0"/>
                <genre-alt value="child_ref" format="fb2.0"/>
                <genre-alt value="child_series" format="fb2.0"/>
                <genre-alt value="child_sports" format="fb2.0"/>
                <genre-alt value="foreign_children" format="flibgolite"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="poetry">
        <root-descr lang="en" genre-title="Poetry, Dramaturgy" detailed="Poetry, Dramaturgy"/>
        <root-descr lang="ru" genre-title="Поэзия, Драматургия" detailed="Поэзия, драматургия"/>
        <root-descr lang="uk" genre-title="Поезія, Драматургія" detailed="Поезія, Драматургія"/>
        <subgenres>
            <subgenre value="poetry">
                <genre-descr lang="en" title="Poetry"/>
                <genre-descr lang="ru" title="Поэзия"/>
                <genre-descr lang="uk" title="Поезія"/>
                <genre-alt value="literature_poetry" format="fb2.0"/>
            </subgenre>
            <subgenre value="dramaturgy">
                <genre-descr lang="en" title="Dramaturgy"/>
                <genre-descr lang="ru" title="Драматургия"/>
                <genre-descr lang="uk" title="Драматургія"/>
                <genre-alt value="performance" format="fb2.0"/>
            </subgenre>
            <subgenre value="humor_verse">
                <genre-descr lang="en" title="Humor Verses"/>
                <genre-descr lang="ru" title="Юмористические Стихи"/>
                <genre-descr lang="uk" title="Гумористичні Вірші"/>
            </subgenre>
            <subgenre value="child_verse">
                <genre-descr lang="en" title="Verses"/>
                <genre-descr lang="ru" title="Детские Стихи"/>
                <genre-descr lang="uk" title="Дитячі Вірші"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="antique">
        <root-descr lang="en" genre-title="Antique" detailed="Antique literature, Myths, Legends"/>
        <root-descr lang="ru" genre-title="Старинное" detailed="Античная литература, мифы, легенды"/>
        <root-descr lang="uk" genre-title="Старовинне" detailed="Антична література, міфи, легенди"/>
        <subgenres>
            <subgenre value="antique_ant">
                <genre-descr lang="en" title="Antique Literature"/>
                <genre-descr lang="ru" title="Античная Литература"/>
                <genre-descr lang="uk" title="Антична література"/>
            </subgenre>
            <subgenre value="antique_european">
                <genre-descr lang="en" title="European Literature"/>
                <genre-descr lang="ru" title="Европейская Старинная Литература"/>
                <genre-descr lang="uk" title="Європейська Старовинна Література"/>
            </subgenre>
            <subgenre value="antique_russian">
                <genre-descr lang="en" title="Antique Russian Literature"/>
                <genre-descr lang="ru" title="Древнерусская Литература"/>
                <genre-descr lang="uk" title="Давньоруська література"/>
            </subgenre>
            <subgenre value="antique_east">
                <genre-descr lang="en" title="Antique East Literature"/>
                <genre-descr lang="ru" title="Древневосточная Литература"/>
                <genre-descr lang="uk" title="Давньосхідна Література"/>
            </subgenre>
            <subgenre value="antique_myths">
                <genre-descr lang="en" title="Myths. Legends. Epos"/>
                <genre-descr lang="ru" title="Мифы. Легенды. Эпос"/>
                <genre-descr lang="uk" title="Міфи. Легенди. Епос"/>
                <genre-alt value="nonfiction_folklor" format="fb2.0"/>
            </subgenre>
            <subgenre value="antique">
                <genre-descr lang="en" title="Other Antique"/>
                <genre-descr lang="ru" title="Старинная Литература: Прочее"/>
                <genre-descr lang="uk" title="Стародавня Література: Інше"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="science">
        <root-descr lang="en" genre-title="Science, Education" detailed="Physics, Philosophy, Psychology... "/>
        <root-descr lang="ru" genre-title="Наука, Образование" detailed="Физика, философия, психология..."/>
        <root-descr lang="uk" genre-title="Наука, Освіта" detailed="Фізика, філософія, психологія..."/>
        <subgenres>
            <subgenre value="sci_history">
                <genre-descr lang="en" title="History"/>
                <genre-descr lang="ru" title="История"/>
                <genre-descr lang="uk" title="Історія"/>
                <genre-alt value="history_africa" format="fb2.0"/>
                <genre-alt value="history_america" format="fb2.0"/>
                <genre-alt value="history_ancient" format="fb2.0"/>
                <genre-alt value="history_asia" format="fb2.0"/>
                <genre-alt value="history_australia" format="fb2.0"/>
                <genre-alt value="history_europe" format="fb2.0"/>
                <genre-alt value="history_study" format="fb2.0"/>
                <genre-alt value="history_jewish" format="fb2.0"/>
                <genre-alt value="history_middle_east" format="fb2.0"/>
                <genre-alt value="histor_military" format="fb2.0"/>
                <genre-alt value="history_military_science" format="fb2.0"/>
                <genre-alt value="history_russia" format="fb2.0"/>
                <genre-alt value="history_usa" format="fb2.0"/>
                <genre-alt value="history_world" format="fb2.0"/>
                <genre-alt value="nonfiction_antropology" format="fb2.0"/>
                <genre-alt value="science_archaeology" format="fb2.0"/>
                <genre-alt value="ref_genealogy" format="fb2.0"/>
                <genre-alt value="science_history_philosophy" format="fb2.0"/>
                <genre-alt value="military_history" format="flibgolite"/>
            </subgenre>
            <subgenre value="sci_psychology">
                <genre-descr lang="en" title="Psychology"/>
                <genre-descr lang="ru" title="Психология"/>
                <genre-descr lang="uk" title="Психологія"/>
                <genre-alt value="health_mental" format="fb2.0"/>
                <genre-alt value="health_psy" format="fb2.0"/>
                <genre-alt value="science_behavioral_sciences" format="fb2.0"/>
                <genre-alt value="science_psy" format="fb2.0"/>
                <genre-alt value="teens_social" format="fb2.0"/>
                <genre-alt value="psy_personal" format="flibgolite"/>
                <genre-alt value="foreign_psychology" format="flibgolite"/>
                <genre-alt value="psy_theraphy" format="flibgolite"/>
                <genre-alt value="psy_generic" format="flibgolite"/>
                <genre-alt value="psy_sex_and_family" format="flibgolite"/>
            </subgenre>
            <subgenre value="sci_culture">
                <genre-descr lang="en" title="Cultural Science"/>
                <genre-descr lang="ru" title="Культурология"/>
                <genre-descr lang="uk" title="Культурологія"/>
            </subgenre>
            <subgenre value="sci_religion">
                <genre-descr lang="en" title="Religious Studies"/>
                <genre-descr lang="ru" title="Религиоведение"/>
                <genre-descr lang="uk" title="Релігієзнавство"/>
            </subgenre>
            <subgenre value="sci_philosophy">
                <genre-descr lang="en" title="Philosophy"/>
                <genre-descr lang="ru" title="Философия"/>
                <genre-descr lang="uk" title="Філософія"/>
                <genre-alt value="nonfiction_philosophy" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_politics">
                <genre-descr lang="en" title="Politics"/>
                <genre-descr lang="ru" title="Политика"/>
                <genre-descr lang="uk" title="Політика"/>
            </subgenre>
            <subgenre value="sci_business">
                <genre-descr lang="en" title="Business literature"/>
                <genre-descr lang="ru" title="Деловая литература"/>
                <genre-descr lang="uk" title="Ділова література"/>
                <genre-alt value="biz_accounting" format="fb2.0"/>
                <genre-alt value="biz_life" format="fb2.0"/>
                <genre-alt value="biz_careers" format="fb2.0"/>
                <genre-alt value="biz_economics" format="fb2.0"/>
                <genre-alt value="biz_finance" format="fb2.0"/>
                <genre-alt value="biz_international" format="fb2.0"/>
                <genre-alt value="biz_professions" format="fb2.0"/>
                <genre-alt value="biz_investing" format="fb2.0"/>
                <genre-alt value="biz_management" format="fb2.0"/>
                <genre-alt value="biz_sales" format="fb2.0"/>
                <genre-alt value="biz_personal_fin" format="fb2.0"/>
                <genre-alt value="biz_ref" format="fb2.0"/>
                <genre-alt value="biz_small_biz" format="fb2.0"/>
                <genre-alt value="professional_finance" format="fb2.0"/>
                <genre-alt value="professional_management" format="fb2.0"/>
                <genre-alt value="ref_edu" format="fb2.0"/>
                <genre-alt value="popular_business" format="flibgolite"/>
                <genre-alt value="management" format="flibgolite"/>
            </subgenre>
            <subgenre value="sci_juris">
                <genre-descr lang="en" title="Jurisprudence"/>
                <genre-descr lang="ru" title="Юриспруденция"/>
                <genre-descr lang="uk" title="Юриспруденція"/>
                <genre-alt value="nonfiction_law" format="fb2.0"/>
                <genre-alt value="professional_law" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_linguistic">
                <genre-descr lang="en" title="Linguistics"/>
                <genre-descr lang="ru" title="Языкознание"/>
                <genre-descr lang="uk" title="Мовазнавство"/>
            </subgenre>
            <subgenre value="sci_medicine">
                <genre-descr lang="en" title="Medicine"/>
                <genre-descr lang="ru" title="Медицина"/>
                <genre-descr lang="uk" title="Медицина"/>
                <genre-alt value="health_aging" format="fb2.0"/>
                <genre-alt value="health_alt_medicine" format="fb2.0"/>
                <genre-alt value="health_cancer" format="fb2.0"/>
                <genre-alt value="professional_medical" format="fb2.0"/>
                <genre-alt value="science_medicine" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_phys">
                <genre-descr lang="en" title="Physics"/>
                <genre-descr lang="ru" title="Физика"/>
                <genre-descr lang="uk" title="Фізика"/>
                <genre-alt value="science_physics" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_math">
                <genre-descr lang="en" title="Mathematics"/>
                <genre-descr lang="ru" title="Математика"/>
                <genre-descr lang="uk" title="Математика"/>
                <genre-alt value="science_math" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_chem">
                <genre-descr lang="en" title="Chemistry"/>
                <genre-descr lang="ru" title="Химия"/>
                <genre-descr lang="uk" title="Хімія"/>
                <genre-alt value="science_chemistry" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_biology">
                <genre-descr lang="en" title="Biology"/>
                <genre-descr lang="ru" title="Биология"/>
                <genre-descr lang="uk" title="Біологія"/>
                <genre-alt value="outdoors_birdwatching" format="fb2.0"/>
                <genre-alt value="outdoors_ecology" format="fb2.0"/>
                <genre-alt value="outdoors_ecosystems" format="fb2.0"/>
                <genre-alt value="outdoors_env" format="fb2.0"/>
                <genre-alt value="outdoors_fauna" format="fb2.0"/>
                <genre-alt value="outdoors_flora" format="fb2.0"/>
                <genre-alt value="outdoors_nature_writing" format="fb2.0"/>
                <genre-alt value="outdoors_ref" format="fb2.0"/>
                <genre-alt value="science_biolog" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_tech">
                <genre-descr lang="en" title="Technical"/>
                <genre-descr lang="ru" title="Технические"/>
                <genre-descr lang="uk" title="Технічні"/>
                <genre-alt value="professional_enginering" format="fb2.0"/>
                <genre-alt value="professional_sci" format="fb2.0"/>
                <genre-alt value="science_technology" format="fb2.0"/>
            </subgenre>
            <subgenre value="science">
                <genre-descr lang="en" title="Misc Science, Education"/>
                <genre-descr lang="ru" title="Научно-образовательная: Прочее"/>
                <genre-descr lang="uk" title="Науково-освітня: Інше"/>
                <genre-alt value="nonfiction_edu" format="fb2.0"/>
                <genre-alt value="nonfiction_gov" format="fb2.0"/>
                <genre-alt value="nonfiction_holidays" format="fb2.0"/>
                <genre-alt value="nonfiction_social_sci" format="fb2.0"/>
                <genre-alt value="nonfiction_ethnology" format="fb2.0"/>
                <genre-alt value="nonfiction_gender" format="fb2.0"/>
                <genre-alt value="nonfiction_gerontology" format="fb2.0"/>
                <genre-alt value="nonfiction_hum_geogr" format="fb2.0"/>
                <genre-alt value="nonfiction_methodology" format="fb2.0"/>
                <genre-alt value="nonfiction_research" format="fb2.0"/>
                <genre-alt value="nonfiction_social_work" format="fb2.0"/>
                <genre-alt value="nonfiction_sociology" format="fb2.0"/>
                <genre-alt value="nonfiction_spec_group" format="fb2.0"/>
                <genre-alt value="nonfiction_stat" format="fb2.0"/>
                <genre-alt value="outdoors_resources" format="fb2.0"/>
                <genre-alt value="professional_edu" format="fb2.0"/>
                <genre-alt value="science_agri" format="fb2.0"/>
                <genre-alt value="science_astronomy" format="fb2.0"/>
                <genre-alt value="science_earth" format="fb2.0"/>
                <genre-alt value="science_edu" format="fb2.0"/>
                <genre-alt value="science_evolution" format="fb2.0"/>
                <genre-alt value="science_measurement" format="fb2.0"/>
                <genre-alt value="science_eco" format="fb2.0"/>
                <genre-alt value="science_ref" format="fb2.0"/>
                <genre-alt value="teens_tech" format="fb2.0"/>
                <genre-alt value="military_special" format="flibgolite"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="computers">
        <root-descr lang="en" genre-title="Computers" detailed="Internet, Programming, Hardware"/>
        <root-descr lang="ru" genre-title="Компьютеры" detailed="Интернет, программирование, железо"/>
        <root-descr lang="uk" genre-title="Комп'ютери" detailed="Інтернет, програмування, залізо"/>
        <subgenres>
            <subgenre value="comp_www">
                <genre-descr lang="en" title="Internet"/>
                <genre-descr lang="ru" title="Интернет"/>
                <genre-descr lang="uk" title="Інтернет"/>
            </subgenre>
            <subgenre value="comp_programming">
                <genre-descr lang="en" title="Programming"/>
                <genre-descr lang="ru" title="Программирование"/>
                <genre-descr lang="uk" title="Програмування"/>
            </subgenre>
            <subgenre value="comp_hard">
                <genre-descr lang="en" title="Hardware"/>
                <genre-descr lang="ru" title="Компьютерное Железо"/>
                <genre-descr lang="uk" title="Комп'ютерне Залізо"/>
                <genre-alt value="comp_hardware" format="fb2.0"/>
            </subgenre>
            <subgenre value="comp_soft">
                <genre-descr lang="en" title="Software"/>
                <genre-descr lang="ru" title="Программы"/>
                <genre-descr lang="uk" title="Програми"/>
                <genre-alt value="comp_software" format="fb2.0"/>
            </subgenre>
            <subgenre value="comp_db">
                <genre-descr lang="en" title="Databases"/>
                <genre-descr lang="ru" title="Базы Данных"/>
                <genre-descr lang="uk" title="Бази Даних"/>
            </subgenre>
            <subgenre value="comp_osnet">
                <genre-descr lang="en" title="OS & Networking"/>
                <genre-descr lang="ru" title="ОС и Сети"/>
                <genre-descr lang="uk" title="ОС та Мережі"/>
                <genre-alt value="comp_microsoft" format="fb2.0"/>
                <genre-alt value="comp_networking" format="fb2.0"/>
                <genre-alt value="comp_os" format="fb2.0"/>
            </subgenre>
            <subgenre value="computers">
                <genre-descr lang="en" title="Computers: Misk"/>
                <genre-descr lang="ru" title="Компьютеры: Прочее"/>
                <genre-descr lang="uk" title="Комп'ютери: Інше"/>
                <genre-alt value="child_computers" format="fb2.0"/>
                <genre-alt value="compusers" format="fb2.0"/>
                <genre-alt value="comp_office" format="fb2.0"/>
                <genre-alt value="comp_cert" format="fb2.0"/>
                <genre-alt value="comp_games" format="fb2.0"/>
                <genre-alt value="comp_sci" format="fb2.0"/>
                <genre-alt value="comp_biz" format="fb2.0"/>
                <genre-alt value="comp_graph" format="fb2.0"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="reference">
        <root-descr lang="en" genre-title="Reference" detailed="Reference, Encyclopedias, Dictionaries"/>
        <root-descr lang="ru" genre-title="Справочники" detailed="Справочники, энциклопедии, словари"/>
        <root-descr lang="uk" genre-title="Довідники" detailed="Довідники, енциклопедії, словники"/>
        <subgenres>
            <subgenre value="ref_encyc">
                <genre-descr lang="en" title="Encyclopedias"/>
                <genre-descr lang="ru" title="Энциклопедии"/>
                <genre-descr lang="uk" title="Енциклопедії"/>
                <genre-alt value="ref_encyclopedia" format="fb2.0"/>
            </subgenre>
            <subgenre value="ref_dict">
                <genre-descr lang="en" title="Dictionaries"/>
                <genre-descr lang="ru" title="Словари"/>
                <genre-descr lang="uk" title="Словники"/>
                <genre-alt value="ref_dict" format="fb2.0"/>
            </subgenre>
            <subgenre value="ref_ref">
                <genre-descr lang="en" title="Reference"/>
                <genre-descr lang="ru" title="Справочники"/>
                <genre-descr lang="uk" title="Довідники"/>
                <genre-alt value="ref_almanacs" format="fb2.0"/>
                <genre-alt value="ref_careers" format="fb2.0"/>
                <genre-alt value="ref_catalogs" format="fb2.0"/>
                <genre-alt value="ref_cons_guides" format="fb2.0"/>
                <genre-alt value="ref_study_guides" format="fb2.0"/>
            </subgenre>
            <subgenre value="ref_guide">
                <genre-descr lang="en" title="Guidebooks"/>
                <genre-descr lang="ru" title="Руководства"/>
                <genre-descr lang="uk" title="Посібники"/>
                <genre-alt value="outdoors_field_guides" format="fb2.0"/>
            </subgenre>
            <subgenre value="reference">
                <genre-descr lang="en" title="Misk References"/>
                <genre-descr lang="ru" title="Справочная Литература: Прочее"/>
                <genre-descr lang="uk" title="Довідкова Література: Інше"/>
                <genre-alt value="nonfiction_ref" format="fb2.0"/>
                <genre-alt value="family_ref" format="fb2.0"/>
                <genre-alt value="references" format="fb2.0"/>
                <genre-alt value="ref_etiquette" format="fb2.0"/>
                <genre-alt value="ref_langs" format="fb2.0"/>
                <genre-alt value="ref_fun" format="fb2.0"/>
                <genre-alt value="ref_books" format="fb2.0"/>
                <genre-alt value="ref_quotations" format="fb2.0"/>
                <genre-alt value="ref_words" format="fb2.0"/>
                <genre-alt value="teens_ref" format="fb2.0"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="nonfiction">
        <root-descr lang="en" genre-title="Nonfiction" detailed="Biography, Memoirs, Publicism"/>
        <root-descr lang="ru" genre-title="Документальное" detailed="Биографии, мемуары, публицистика"/>
        <root-descr lang="uk" genre-title="Документальне" detailed="Біографії, мемуари, публіцистика"/>
        <subgenres>
            <subgenre value="nonf_biography">
                <genre-descr lang="en" title="Biography & Memoirs"/>
                <genre-descr lang="ru" title="Биографии и Мемуары"/>
                <genre-descr lang="uk" title="Біографії та Мемуари"/>
                <genre-alt value="people" format="fb2.0"/>
                <genre-alt value="biography" format="fb2.0"/>
                <genre-alt value="biogr_arts" format="fb2.0"/>
                <genre-alt value="biogr_ethnic" format="fb2.0"/>
                <genre-alt value="biogr_family" format="fb2.0"/>
                <genre-alt value="biogr_historical" format="fb2.0"/>
                <genre-alt value="biogr_leaders" format="fb2.0"/>
                <genre-alt value="biogr_professionals" format="fb2.0"/>
                <genre-alt value="biogr_sports" format="fb2.0"/>
                <genre-alt value="biogr_travel" format="fb2.0"/>
                <genre-alt value="biz_beogr" format="fb2.0"/>
                <genre-alt value="gay_biogr" format="fb2.0"/>
                <genre-alt value="history_gay" format="fb2.0"/>
                <genre-alt value="literature_letters" format="fb2.0"/>
                <genre-alt value="teens_beogr" format="fb2.0"/>
            </subgenre>
            <subgenre value="nonf_publicism">
                <genre-descr lang="en" title="Publicism"/>
                <genre-descr lang="ru" title="Публицистика"/>
                <genre-descr lang="uk" title="Публіцистика"/>
                <genre-alt value="foreign_publicism" format="flibgolite"/>
            </subgenre>
            <subgenre value="nonf_criticism">
                <genre-descr lang="en" title="Criticism"/>
                <genre-descr lang="ru" title="Критика"/>
                <genre-descr lang="uk" title="Критика"/>
            </subgenre>
            <subgenre value="nonfiction">
                <genre-descr lang="en" title="Misk Nonfiction"/>
                <genre-descr lang="ru" title="Документальное: Прочее"/>
                <genre-descr lang="uk" title="Документальне: Інше"/>
                <genre-alt value="gay_nonfiction" format="fb2.0"/>
                <genre-alt value="nonfiction_avto" format="fb2.0"/>
                <genre-alt value="nonfiction_crime" format="fb2.0"/>
                <genre-alt value="nonfiction_events" format="fb2.0"/>
                <genre-alt value="nonfiction_politics" format="fb2.0"/>
                <genre-alt value="nonfiction_traditions" format="fb2.0"/>
                <genre-alt value="nonfiction_demography" format="fb2.0"/>
                <genre-alt value="nonfiction_racism" format="fb2.0"/>
                <genre-alt value="nonfiction_emigration" format="fb2.0"/>
                <genre-alt value="nonfiction_philantropy" format="fb2.0"/>
                <genre-alt value="nonfiction_transportation" format="fb2.0"/>
                <genre-alt value="nonfiction_true_accounts" format="fb2.0"/>
                <genre-alt value="nonfiction_urban" format="fb2.0"/>
                <genre-alt value="nonfiction_women" format="fb2.0"/>
                <genre-alt value="outdoors_conservation" format="fb2.0"/>
            </subgenre>
            <subgenre value="design">
                <genre-descr lang="en" title="Art, Design"/>
                <genre-descr lang="ru" title="Искусство, Дизайн"/>
                <genre-descr lang="uk" title="Мистецтво, Дизайн"/>
                <genre-alt value="architecture" format="fb2.0"/>
                <genre-alt value="art" format="fb2.0"/>
                <genre-alt value="art_instr" format="fb2.0"/>
                <genre-alt value="artists" format="fb2.0"/>
                <genre-alt value="fashion" format="fb2.0"/>
                <genre-alt value="graph_design" format="fb2.0"/>
                <genre-alt value="photography" format="fb2.0"/>
                <genre-alt value="music" format="flibgolite"/>
                <genre-alt value="notes" format="flibgolite"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="religion">
        <root-descr lang="en" genre-title="Religion" detailed="Religion, Esoterics"/>
        <root-descr lang="ru" genre-title="Религия" detailed="Религия, эзотерика"/>
        <root-descr lang="uk" genre-title="Релігія" detailed="Релігія, езотерика"/>
        <subgenres>
            <subgenre value="religion_rel">
                <genre-descr lang="en" title="Religion"/>
                <genre-descr lang="ru" title="Религия"/>
                <genre-descr lang="uk" title="Релігія"/>
                <genre-alt value="child_religion" format="fb2.0"/>
                <genre-alt value="chris_bibles" format="fb2.0"/>
                <genre-alt value="chris_pravoslavie" format="fb2.0"/>
                <genre-alt value="chris_catholicism" format="fb2.0"/>
                <genre-alt value="chris_living" format="fb2.0"/>
                <genre-alt value="chris_history" format="fb2.0"/>
                <genre-alt value="chris_clergy" format="fb2.0"/>
                <genre-alt value="chris_edu" format="fb2.0"/>
                <genre-alt value="chris_evangelism" format="fb2.0"/>
                <genre-alt value="chris_fiction" format="fb2.0"/>
                <genre-alt value="chris_holidays" format="fb2.0"/>
                <genre-alt value="chris_jesus" format="fb2.0"/>
                <genre-alt value="chris_mormonism" format="fb2.0"/>
                <genre-alt value="chris_orthodoxy" format="fb2.0"/>
                <genre-alt value="outdoors_conservation" format="fb2.0"/>
                <genre-alt value="chris_protestantism" format="fb2.0"/>
                <genre-alt value="chris_ref" format="fb2.0"/>
                <genre-alt value="chris_theology" format="fb2.0"/>
                <genre-alt value="chris_devotion" format="fb2.0"/>
                <genre-alt value="literature_religion" format="fb2.0"/>
                <genre-alt value="religion" format="fb2.0"/>
                <genre-alt value="religion_bibles" format="fb2.0"/>
                <genre-alt value="Christianity" format="fb2.0"/>
                <genre-alt value="religion_fiction" format="fb2.0"/>
                <genre-alt value="religion_new_age" format="fb2.0"/>
                <genre-alt value="religion_religious_studies" format="fb2.0"/>
                <genre-alt value="romance_religion" format="fb2.0"/>
                <genre-alt value="teens_religion" format="fb2.0"/>
                <genre-alt value="religion_orthodoxy" format="flibgolite"/>
                <genre-alt value="religion_christianity" format="flibgolite"/>
            </subgenre>
            <subgenre value="religion_esoterics">
                <genre-descr lang="en" title="Esoterics"/>
                <genre-descr lang="ru" title="Эзотерика"/>
                <genre-descr lang="uk" title="Езотерика"/>
                <genre-alt value="religion_occult" format="fb2.0"/>
                <genre-alt value="religion_spirituality" format="fb2.0"/>
            </subgenre>
            <subgenre value="religion_self">
                <genre-descr lang="en" title="Self-perfection"/>
                <genre-descr lang="ru" title="Самосовершенствование"/>
                <genre-descr lang="uk" title="Самовдосконалення"/>
            </subgenre>
            <subgenre value="religion">
                <genre-descr lang="en" title="Religion: Other"/>
                <genre-descr lang="ru" title="Религия и духовность: Прочее"/>
                <genre-descr lang="uk" title="Релігія та духовність: Інше"/>
                <genre-alt value="religion_east" format="fb2.0"/>
                <genre-alt value="religion_buddhism" format="fb2.0"/>
                <genre-alt value="religion_earth" format="fb2.0"/>
                <genre-alt value="religion_hinduism" format="fb2.0"/>
                <genre-alt value="religion_islam" format="fb2.0"/>
                <genre-alt value="religion_judaism" format="fb2.0"/>
                <genre-alt value="religion_other" format="fb2.0"/>
            </subgenre>
            <subgenre value="sci_religion">
                <genre-descr lang="en" title="Religious Studies"/>
                <genre-descr lang="ru" title="Религиоведение"/>
                <genre-descr lang="uk" title="Релігієзнавство"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="humor">
        <root-descr lang="en" genre-title="Humor" detailed="Prose, Verses, Anecdote"/>
        <root-descr lang="ru" genre-title="Юмор" detailed="Проза, стихи, анекдоты"/>
        <root-descr lang="uk" genre-title="Гумор" detailed="Проза, вірші, анекдоти"/>
        <subgenres>
            <subgenre value="humor_anecdote">
                <genre-descr lang="en" title="Anecdote"/>
                <genre-descr lang="ru" title="Анекдоты"/>
                <genre-descr lang="uk" title="Анекдоти"/>
            </subgenre>
            <subgenre value="humor_prose">
                <genre-descr lang="en" title="Humor Prose"/>
                <genre-descr lang="ru" title="Юмористическая Проза"/>
                <genre-descr lang="uk" title="Гумористична Проза"/>
            </subgenre>
            <subgenre value="humor_verse">
                <genre-descr lang="en" title="Humor Verses"/>
                <genre-descr lang="ru" title="Юмористические Стихи"/>
                <genre-descr lang="uk" title="Гумористичні Вірші"/>
            </subgenre>
            <subgenre value="humor">
                <genre-descr lang="en" title="Misc Humor"/>
                <genre-descr lang="ru" title="Юмор: Прочее"/>
                <genre-descr lang="uk" title="Гумор: Інше"/>
                <genre-alt value="family_humor" format="fb2.0"/>
            </subgenre>
            <subgenre value="humor_satire">
                <genre-descr lang="en" title="Satire"/>
                <genre-descr lang="ru" title="Сатира"/>
                <genre-descr lang="uk" title="Сатира"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="home">
        <root-descr lang="en" genre-title="Home, Family" detailed="Cooking, Pets, Hobby"/>
        <root-descr lang="ru" genre-title="Дом, Семья" detailed="Кулинария, домашние животные, хобби"/>
        <root-descr lang="uk" genre-title="Дім, Сім'я" detailed="Кулінарія, домашні тварини, хобі"/>
        <subgenres>
            <subgenre value="home_cooking">
                <genre-descr lang="en" title="Cooking"/>
                <genre-descr lang="ru" title="Кулинария"/>
                <genre-descr lang="uk" title="Кулінарія"/>
                <genre-alt value="cooking" format="fb2.0"/>
                <genre-alt value="cook_baking" format="fb2.0"/>
                <genre-alt value="cook_can" format="fb2.0"/>
                <genre-alt value="cook_art" format="fb2.0"/>
                <genre-alt value="cook_drink" format="fb2.0"/>
                <genre-alt value="cook_gastronomy" format="fb2.0"/>
                <genre-alt value="cook_meals" format="fb2.0"/>
                <genre-alt value="cook_natura" format="fb2.0"/>
                <genre-alt value="cook_outdoor" format="fb2.0"/>
                <genre-alt value="cook_pro" format="fb2.0"/>
                <genre-alt value="cook_quick" format="fb2.0"/>
                <genre-alt value="cook_ref" format="fb2.0"/>
                <genre-alt value="cook_regional" format="fb2.0"/>
                <genre-alt value="cook_appliances" format="fb2.0"/>
                <genre-alt value="cook_diet" format="fb2.0"/>
                <genre-alt value="cook_spec" format="fb2.0"/>
                <genre-alt value="cook_veget" format="fb2.0"/>
                <genre-alt value="health_diets" format="fb2.0"/>
            </subgenre>
            <subgenre value="home_pets">
                <genre-descr lang="en" title="Pets"/>
                <genre-descr lang="ru" title="Домашние Животные"/>
                <genre-descr lang="uk" title="Домашні Тварини"/>
            </subgenre>
            <subgenre value="home_crafts">
                <genre-descr lang="en" title="Hobbies & Crafts"/>
                <genre-descr lang="ru" title="Хобби, Ремесла"/>
                <genre-descr lang="uk" title="Хобі, Ремесла"/>
                <genre-alt value="home_collect" format="fb2.0"/>
                <genre-alt value="outdoors_hiking" format="fb2.0"/>
                <genre-alt value="outdoors_hunt_fish" format="fb2.0"/>
            </subgenre>
            <subgenre value="home_entertain">
                <genre-descr lang="en" title="Entertaining"/>
                <genre-descr lang="ru" title="Развлечения"/>
                <genre-descr lang="uk" title="Розваги"/>
                <genre-alt value="entertainment" format="fb2.0"/>
                <genre-alt value="entert_comics" format="fb2.0"/>
                <genre-alt value="entert_games" format="fb2.0"/>
                <genre-alt value="entert_humor" format="fb2.0"/>
                <genre-alt value="entert_movies" format="fb2.0"/>
                <genre-alt value="entert_music" format="fb2.0"/>
                <genre-alt value="nonfiction_pop_culture" format="fb2.0"/>
                <genre-alt value="entert_radio" format="fb2.0"/>
                <genre-alt value="entert_tv" format="fb2.0"/>
            </subgenre>
            <subgenre value="home_health">
                <genre-descr lang="en" title="Health"/>
                <genre-descr lang="ru" title="Здоровье"/>
                <genre-descr lang="uk" title="Здоров'я"/>
                <genre-alt value="health" format="fb2.0"/>
                <genre-alt value="health_beauty" format="fb2.0"/>
                <genre-alt value="family_health" format="fb2.0"/>
                <genre-alt value="family_fertility" format="fb2.0"/>
                <genre-alt value="family_parenting" format="fb2.0"/>
                <genre-alt value="family_pregnancy" format="fb2.0"/>
                <genre-alt value="family_special_needs" format="fb2.0"/>
                <genre-alt value="health_death" format="fb2.0"/>
                <genre-alt value="health_dideases" format="fb2.0"/>
                <genre-alt value="health_fitness" format="fb2.0"/>
                <genre-alt value="health_men" format="fb2.0"/>
                <genre-alt value="health_nutrition" format="fb2.0"/>
                <genre-alt value="health_personal" format="fb2.0"/>
                <genre-alt value="health_recovery" format="fb2.0"/>
                <genre-alt value="health_ref" format="fb2.0"/>
                <genre-alt value="health_first_aid" format="fb2.0"/>
                <genre-alt value="health_self_help" format="fb2.0"/>
                <genre-alt value="health_women" format="fb2.0"/>
            </subgenre>
            <subgenre value="home_garden">
                <genre-descr lang="en" title="Garden"/>
                <genre-descr lang="ru" title="Сад и Огород"/>
                <genre-descr lang="uk" title="Сад і Город"/>
            </subgenre>
            <subgenre value="home_diy">
                <genre-descr lang="en" title="Do it yourself"/>
                <genre-descr lang="ru" title="Сделай Сам"/>
                <genre-descr lang="uk" title="Зроби Сам"/>
                <genre-alt value="home_expert" format="fb2.0"/>
                <genre-alt value="home_design" format="fb2.0"/>
                <genre-alt value="home_howto" format="fb2.0"/>
                <genre-alt value="home_interior_design" format="fb2.0"/>
            </subgenre>
            <subgenre value="home_sport">
                <genre-descr lang="en" title="Sports"/>
                <genre-descr lang="ru" title="Спорт"/>
                <genre-descr lang="uk" title="Спорт"/>
                <genre-alt value="literature_sports" format="fb2.0"/>
                <genre-alt value="outdoors_outdoor_recreation" format="fb2.0"/>
                <genre-alt value="outdoors_survive" format="fb2.0"/>
                <genre-alt value="sport" format="fb2.0"/>
                <genre-alt value="teens_health" format="fb2.0"/>
                <genre-alt value="teens_school_sports" format="fb2.0"/>
            </subgenre>
            <subgenre value="home_sex">
                <genre-descr lang="en" title="Erotica, Sex"/>
                <genre-descr lang="ru" title="Эротика, Секс"/>
                <genre-descr lang="uk" title="Еротика, Секс"/>
                <genre-alt value="health_sex" format="fb2.0"/>
                <genre-alt value="nonfiction_pornography" format="fb2.0"/>
            </subgenre>
            <subgenre value="home">
                <genre-descr lang="en" title="Home: Other"/>
                <genre-descr lang="ru" title="Дом и Семья: Прочее"/>
                <genre-descr lang="uk" title="Дім та Сім'я: Інше"/>
                <genre-alt value="gay_parenting" format="fb2.0"/>
                <genre-alt value="home_cottage" format="fb2.0"/>
                <genre-alt value="home_weddings" format="fb2.0"/>
                <genre-alt value="family" format="fb2.0"/>
                <genre-alt value="family_adoption" format="fb2.0"/>
                <genre-alt value="family_aging_parents" format="fb2.0"/>
                <genre-alt value="family_edu" format="fb2.0"/>
                <genre-alt value="family_activities" format="fb2.0"/>
                <genre-alt value="family_relations" format="fb2.0"/>
                <genre-alt value="family_lit_guide" format="fb2.0"/>
                <genre-alt value="women_divorce" format="fb2.0"/>
                <genre-alt value="women_domestic" format="fb2.0"/>
                <genre-alt value="women_child" format="fb2.0"/>
                <genre-alt value="women_single" format="fb2.0"/>
            </subgenre>
        </subgenres>
    </genre>
</fbgenrestransfer>
`

const ALT_GENRES_XML = `
<!-- borrowed from SOPDS -->
<?xml version="1.0" encoding="utf-8"?>
<fbgenrestransfer xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="http://alexs.ru/fb2/genrestable/GT.xsd">
    <genre value="economics_ref">
        <root-descr lang="ru" genre-title="Деловая литература" detailed=""/>
        <subgenres>
            <subgenre value="economics_ref">
                <genre-descr lang="ru" title="Деловая литература"/>
            </subgenre>
            <subgenre value="popular_business">
                <genre-descr lang="ru" title="Карьера, кадры"/>
            </subgenre>
            <subgenre value="org_behavior">
                <genre-descr lang="ru" title="Маркетинг, PR"/>
            </subgenre>
            <subgenre value="banking">
                <genre-descr lang="ru" title="Финансы"/>
            </subgenre>
            <subgenre value="economics">
                <genre-descr lang="ru" title="Экономика"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="detective">
        <root-descr lang="ru" genre-title="Детективы и Триллеры" detailed=""/>
        <subgenres>
            <subgenre value="detective">
                <genre-descr lang="ru" title="Детективы"/>
            </subgenre>
            <subgenre value="det_action">
                <genre-descr lang="ru" title="Боевик"/>
            </subgenre>
            <subgenre value="det_irony">
                <genre-descr lang="ru" title="Иронический детектив, дамский детективный роман"/>
            </subgenre>
            <subgenre value="det_history">
                <genre-descr lang="ru" title="Исторический детектив"/>
            </subgenre>
            <subgenre value="det_classic">
                <genre-descr lang="ru" title="Классический детектив"/>
            </subgenre>
            <subgenre value="det_crime">
                <genre-descr lang="ru" title="Криминальный детектив"/>
            </subgenre>
            <subgenre value="det_hard">
                <genre-descr lang="ru" title="Крутой детектив"/>
            </subgenre>
            <subgenre value="det_political">
                <genre-descr lang="ru" title="Политический детектив"/>
            </subgenre>
            <subgenre value="det_police">
                <genre-descr lang="ru" title="Полицейский детектив"/>
            </subgenre>
            <subgenre value="det_maniac">
                <genre-descr lang="ru" title="Про маньяков"/>
            </subgenre>
            <subgenre value="det_su">
                <genre-descr lang="ru" title="Советский детектив"/>
            </subgenre>
            <subgenre value="thriller">
                <genre-descr lang="ru" title="Триллер"/>
            </subgenre>
            <subgenre value="det_espionage">
                <genre-descr lang="ru" title="Шпионский детектив"/>
            </subgenre>
            <subgenre value="det_action">
                <genre-descr lang="ru" title="Боевик"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="nonfiction">
        <root-descr lang="ru" genre-title="Документальная литература" detailed=""/>
        <subgenres>
            <subgenre value="nonfiction">
                <genre-descr lang="ru" title="Документальная литература"/>
            </subgenre>
            <subgenre value="nonf_biography">
                <genre-descr lang="ru" title="Биографии и Мемуары"/>
            </subgenre>
            <subgenre value="nonf_military">
                <genre-descr lang="ru" title="Военная документалистика и аналитика"/>
            </subgenre>
            <subgenre value="military_special">
                <genre-descr lang="ru" title="Военное дело"/>
            </subgenre>
            <subgenre value="travel_notes">
                <genre-descr lang="ru" title="География, путевые заметки"/>
            </subgenre>
            <subgenre value="nonf_publicism">
                <genre-descr lang="ru" title="Публицистика"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="home">
        <root-descr lang="ru" genre-title="Дом и семья" detailed=""/>
        <subgenres>
            <subgenre value="home">
                <genre-descr lang="ru" title="Домоводство"/>
            </subgenre>
            <subgenre value="auto_regulations">
                <genre-descr lang="ru" title="Автомобили и ПДД"/>
            </subgenre>
            <subgenre value="home_sport">
                <genre-descr lang="ru" title="Боевые искусства, спорт"/>
            </subgenre>
            <subgenre value="home_pets">
                <genre-descr lang="ru" title="Домашние животные"/>
            </subgenre>
            <subgenre value="home_health">
                <genre-descr lang="ru" title="Здоровье"/>
            </subgenre>
            <subgenre value="home_collecting">
                <genre-descr lang="ru" title="Коллекционирование"/>
            </subgenre>
            <subgenre value="home_cooking">
                <genre-descr lang="ru" title="Кулинария"/>
            </subgenre>
            <subgenre value="sci_pedagogy">
                <genre-descr lang="ru" title="Педагогика, воспитание детей, литература для родителей"/>
            </subgenre>
            <subgenre value="home_entertain">
                <genre-descr lang="ru" title="Развлечения"/>
            </subgenre>
            <subgenre value="home_garden">
                <genre-descr lang="ru" title="Сад и огород"/>
            </subgenre>
            <subgenre value="home_diy">
                <genre-descr lang="ru" title="Сделай сам"/>
            </subgenre>
            <subgenre value="family">
                <genre-descr lang="ru" title="Семейные отношения"/>
            </subgenre>
            <subgenre value="home_sex">
                <genre-descr lang="ru" title="Семейные отношения, секс"/>
            </subgenre>
            <subgenre value="home_crafts">
                <genre-descr lang="ru" title="Хобби и ремесла"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="drama">
        <root-descr lang="ru" genre-title="Драматургия" detailed=""/>
        <subgenres>
            <subgenre value="drama">
                <genre-descr lang="ru" title="Драма"/>
            </subgenre>
            <subgenre value="drama_antique">
                <genre-descr lang="ru" title="Античная драма"/>
            </subgenre>
            <subgenre value="dramaturgy">
                <genre-descr lang="ru" title="Драматургия"/>
            </subgenre>
            <subgenre value="comedy">
                <genre-descr lang="ru" title="Комедия"/>
            </subgenre>
            <subgenre value="vaudeville">
                <genre-descr lang="ru" title="Мистерия, буффонада, водевиль"/>
            </subgenre>
            <subgenre value="screenplays">
                <genre-descr lang="ru" title="Сценарий"/>
            </subgenre>
            <subgenre value="tragedy">
                <genre-descr lang="ru" title="Трагедия"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="design">
        <root-descr lang="ru" genre-title="Искусство, Искусствоведение, Дизайн" detailed=""/>
        <subgenres>
            <subgenre value="painting">
                <genre-descr lang="ru" title="Живопись, альбомы, иллюстрированные каталоги"/>
            </subgenre>
            <subgenre value="design">
                <genre-descr lang="ru" title="Искусство и Дизайн"/>
            </subgenre>
            <subgenre value="art_criticism">
                <genre-descr lang="ru" title="Искусствоведение"/>
            </subgenre>
            <subgenre value="cine">
                <genre-descr lang="ru" title="Кино"/>
            </subgenre>
            <subgenre value="nonf_criticism">
                <genre-descr lang="ru" title="Критика"/>
            </subgenre>
            <subgenre value="sci_culture">
                <genre-descr lang="ru" title="Культурология"/>
            </subgenre>
            <subgenre value="art_world_culture">
                <genre-descr lang="ru" title="Мировая художественная культура"/>
            </subgenre>
            <subgenre value="music">
                <genre-descr lang="ru" title="Музыка"/>
            </subgenre>
            <subgenre value="notes">
                <genre-descr lang="ru" title="Партитуры"/>
            </subgenre>
            <subgenre value="architecture_book">
                <genre-descr lang="ru" title="Скульптура и архитектура"/>
            </subgenre>
            <subgenre value="theatre">
                <genre-descr lang="ru" title="Театр"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="computers">
        <root-descr lang="ru" genre-title="Компьютеры и Интернет" detailed=""/>
        <subgenres>
            <subgenre value="computers">
                <genre-descr lang="ru" title="Компьютерная, околокомпьютерная литература"/>
            </subgenre>
            <subgenre value="comp_hard">
                <genre-descr lang="ru" title="Компьютерное железо"/>
            </subgenre>
            <subgenre value="comp_www">
                <genre-descr lang="ru" title="ОС и Сети, Интернет"/>
            </subgenre>
            <subgenre value="comp_db">
                <genre-descr lang="ru" title="Программирование, программы, базы данных"/>
            </subgenre>
            <subgenre value="tbg_computers">
                <genre-descr lang="ru" title="Учебные пособия, самоучители"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="children">
        <root-descr lang="ru" genre-title="Литература для детей" detailed=""/>
        <subgenres>
            <subgenre value="children">
                <genre-descr lang="ru" title="Детская литература"/>
            </subgenre>
            <subgenre value="child_education">
                <genre-descr lang="ru" title="Детская образовательная литература"/>
            </subgenre>
            <subgenre value="child_det">
                <genre-descr lang="ru" title="Детская остросюжетная литература"/>
            </subgenre>
            <subgenre value="foreign_children">
                <genre-descr lang="ru" title="Зарубежная литература для детей"/>
            </subgenre>
            <subgenre value="prose_game">
                <genre-descr lang="ru" title="Игры, упражнения для детей"/>
            </subgenre>
            <subgenre value="child_classical">
                <genre-descr lang="ru" title="Классическая детская литература"/>
            </subgenre>
            <subgenre value="child_prose">
                <genre-descr lang="ru" title="Проза для детей"/>
            </subgenre>
            <subgenre value="child_tale_rus">
                <genre-descr lang="ru" title="Русские сказки"/>
            </subgenre>
            <subgenre value="child_tale">
                <genre-descr lang="ru" title="Сказки народов мира"/>
            </subgenre>
            <subgenre value="child_verse">
                <genre-descr lang="ru" title="Фантастика для детей"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="love">
        <root-descr lang="ru" genre-title="Любовные романы" detailed=""/>
        <subgenres>
            <subgenre value="love_history">
                <genre-descr lang="ru" title="Исторические любовные романы"/>
            </subgenre>
            <subgenre value="love_short">
                <genre-descr lang="ru" title="Короткие любовные романы"/>
            </subgenre>
            <subgenre value="love_sf">
                <genre-descr lang="ru" title="Любовное фэнтези, любовно-фантастические романы"/>
            </subgenre>
            <subgenre value="love">
                <genre-descr lang="ru" title="Любовные романы"/>
            </subgenre>
            <subgenre value="love_detective">
                <genre-descr lang="ru" title="Остросюжетные любовные романы"/>
            </subgenre>
            <subgenre value="love_hard">
                <genre-descr lang="ru" title="Порно"/>
            </subgenre>
            <subgenre value="love_contemporary">
                <genre-descr lang="ru" title="Современные любовные романы"/>
            </subgenre>
            <subgenre value="love_erotica">
                <genre-descr lang="ru" title="Эротическая литература"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="science">
        <root-descr lang="ru" genre-title="Наука, Образование" detailed=""/>
        <subgenres>
            <subgenre value="sci_medicine_alternative">
                <genre-descr lang="ru" title="Альтернативная медицина"/>
            </subgenre>
            <subgenre value="sci_theories">
                <genre-descr lang="ru" title="Альтернативные науки и научные теории"/>
            </subgenre>
            <subgenre value="sci_cosmos">
                <genre-descr lang="ru" title="Астрономия и Космос"/>
            </subgenre>
            <subgenre value="sci_biology">
                <genre-descr lang="ru" title="Биология, биофизика, биохимия"/>
            </subgenre>
            <subgenre value="sci_botany">
                <genre-descr lang="ru" title="Ботаника"/>
            </subgenre>
            <subgenre value="sci_veterinary">
                <genre-descr lang="ru" title="Ветеринария"/>
            </subgenre>
            <subgenre value="military_history">
                <genre-descr lang="ru" title="Военная история"/>
            </subgenre>
            <subgenre value="sci_oriental">
                <genre-descr lang="ru" title="Востоковедение"/>
            </subgenre>
            <subgenre value="sci_geo">
                <genre-descr lang="ru" title="Геология и география"/>
            </subgenre>
            <subgenre value="sci_state">
                <genre-descr lang="ru" title="Государство и право"/>
            </subgenre>
            <subgenre value="sci_popular">
                <genre-descr lang="ru" title="Зарубежная образовательная литература, зарубежная прикладная, научно-популярная литература"/>
            </subgenre>
            <subgenre value="sci_zoo">
                <genre-descr lang="ru" title="Зоология"/>
            </subgenre>
            <subgenre value="sci_history">
                <genre-descr lang="ru" title="История"/>
            </subgenre>
            <subgenre value="sci_philology">
                <genre-descr lang="ru" title="Литературоведение"/>
            </subgenre>
            <subgenre value="sci_math">
                <genre-descr lang="ru" title="Математика"/>
            </subgenre>
            <subgenre value="science">
                <genre-descr lang="ru" title="Научная литература"/>
            </subgenre>
            <subgenre value="sci_social_studies">
                <genre-descr lang="ru" title="Обществознание, социология"/>
            </subgenre>
            <subgenre value="sci_politics">
                <genre-descr lang="ru" title="Политика"/>
            </subgenre>
            <subgenre value="sci_psychology">
                <genre-descr lang="ru" title="Психология и психотерапия"/>
            </subgenre>
            <subgenre value="sci_phys">
                <genre-descr lang="ru" title="Физика"/>
            </subgenre>
            <subgenre value="sci_philosophy">
                <genre-descr lang="ru" title="Философия"/>
            </subgenre>
            <subgenre value="sci_chem">
                <genre-descr lang="ru" title="Химия"/>
            </subgenre>
            <subgenre value="sci_ecology">
                <genre-descr lang="ru" title="Экология"/>
            </subgenre>
            <subgenre value="sci_economy">
                <genre-descr lang="ru" title="Экономика"/>
            </subgenre>
            <subgenre value="sci_juris">
                <genre-descr lang="ru" title="Юриспруденция"/>
            </subgenre>
            <subgenre value="sci_linguistic">
                <genre-descr lang="ru" title="Языкознание, иностранные языки"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="poetry">
        <root-descr lang="ru" genre-title="Поэзия" detailed=""/>
        <subgenres>
            <subgenre value="palindromes">
                <genre-descr lang="ru" title="Визуальная и экспериментальная поэзия, верлибры, палиндромы"/>
            </subgenre>
            <subgenre value="poetry_for_classical">
                <genre-descr lang="ru" title="Классическая зарубежная поэзия"/>
            </subgenre>
            <subgenre value="poetry_classical">
                <genre-descr lang="ru" title="Классическая поэзия"/>
            </subgenre>
            <subgenre value="poetry_rus_classical">
                <genre-descr lang="ru" title="Классическая русская поэзия"/>
            </subgenre>
            <subgenre value="lyrics">
                <genre-descr lang="ru" title="Лирика"/>
            </subgenre>
            <subgenre value="song_poetry">
                <genre-descr lang="ru" title="Песенная поэзия"/>
            </subgenre>
            <subgenre value="poetry">
                <genre-descr lang="ru" title="Поэзия"/>
            </subgenre>
            <subgenre value="poetry_east">
                <genre-descr lang="ru" title="Поэзия Востока"/>
            </subgenre>
            <subgenre value="poem">
                <genre-descr lang="ru" title="Поэма, эпическая поэзия"/>
            </subgenre>
            <subgenre value="poetry_for_modern">
                <genre-descr lang="ru" title="Современная зарубежная поэзия"/>
            </subgenre>
            <subgenre value="poetry_modern">
                <genre-descr lang="ru" title="Современная поэзия"/>
            </subgenre>
            <subgenre value="poetry_rus_modern">
                <genre-descr lang="ru" title="Современная русская поэзия"/>
            </subgenre>
            <subgenre value="humor_verse">
                <genre-descr lang="ru" title="Юмористические стихи, басни"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="adventure">
        <root-descr lang="ru" genre-title="Приключения" detailed=""/>
        <subgenres>
            <subgenre value="adv_story">
                <genre-descr lang="ru" title="Авантюрный роман"/>
            </subgenre>
            <subgenre value="adv_indian">
                <genre-descr lang="ru" title="Вестерн, про индейцев"/>
            </subgenre>
            <subgenre value="adv_history">
                <genre-descr lang="ru" title="Исторические приключения"/>
            </subgenre>
            <subgenre value="adv_maritime">
                <genre-descr lang="ru" title="Морские приключения"/>
            </subgenre>
            <subgenre value="adventure">
                <genre-descr lang="ru" title="Приключения"/>
            </subgenre>
            <subgenre value="adv_modern">
                <genre-descr lang="ru" title="Приключения в современном мире"/>
            </subgenre>
            <subgenre value="child_adv">
                <genre-descr lang="ru" title="Приключения для детей и подростков"/>
            </subgenre>
            <subgenre value="adv_animal">
                <genre-descr lang="ru" title="Природа и животные"/>
            </subgenre>
            <subgenre value="adv_geo">
                <genre-descr lang="ru" title="Путешествия и география"/>
            </subgenre>
            <subgenre value="tale_chivalry">
                <genre-descr lang="ru" title="Рыцарский роман"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="prose">
        <root-descr lang="ru" genre-title="Проза" detailed=""/>
        <subgenres>
            <subgenre value="aphorisms">
                <genre-descr lang="ru" title="Афоризмы, цитаты"/>
            </subgenre>
            <subgenre value="gothic_novel">
                <genre-descr lang="ru" title="Готический роман"/>
            </subgenre>
            <subgenre value="foreign_prose">
                <genre-descr lang="ru" title="Зарубежная классическая проза"/>
            </subgenre>
            <subgenre value="prose_history">
                <genre-descr lang="ru" title="Историческая проза"/>
            </subgenre>
            <subgenre value="prose_classic">
                <genre-descr lang="ru" title="Классическая проза"/>
            </subgenre>
            <subgenre value="literature_18">
                <genre-descr lang="ru" title="Классическая проза XVII-XVIII веков"/>
            </subgenre>
            <subgenre value="literature_19">
                <genre-descr lang="ru" title="Классическая проза ХIX века"/>
            </subgenre>
            <subgenre value="literature_20">
                <genre-descr lang="ru" title="Классическая проза ХX века"/>
            </subgenre>
            <subgenre value="prose_counter">
                <genre-descr lang="ru" title="Контркультура"/>
            </subgenre>
            <subgenre value="prose_magic">
                <genre-descr lang="ru" title="Магический реализм"/>
            </subgenre>
            <subgenre value="story">
                <genre-descr lang="ru" title="Малые литературные формы прозы: рассказы, эссе, новеллы, феерия"/>
            </subgenre>
            <subgenre value="prose">
                <genre-descr lang="ru" title="Проза"/>
            </subgenre>
            <subgenre value="prose_military">
                <genre-descr lang="ru" title="Проза о войне"/>
            </subgenre>
            <subgenre value="great_story">
                <genre-descr lang="ru" title="Роман, повесть"/>
            </subgenre>
            <subgenre value="prose_rus_classic">
                <genre-descr lang="ru" title="Русская классическая проза"/>
            </subgenre>
            <subgenre value="prose_su_classics">
                <genre-descr lang="ru" title="Советская классическая проза"/>
            </subgenre>
            <subgenre value="prose_contemporary">
                <genre-descr lang="ru" title="Современная русская и зарубежная проза"/>
            </subgenre>
            <subgenre value="foreign_antique">
                <genre-descr lang="ru" title="Средневековая классическая проза"/>
            </subgenre>
            <subgenre value="prose_abs">
                <genre-descr lang="ru" title="Фантасмагория, абсурдистская проза"/>
            </subgenre>
            <subgenre value="prose_neformatny">
                <genre-descr lang="ru" title="Экспериментальная, неформатная проза"/>
            </subgenre>
            <subgenre value="epistolary_fiction">
                <genre-descr lang="ru" title="Эпистолярная проза"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="other">
        <root-descr lang="ru" genre-title="Прочее" detailed=""/>
        <subgenres>
            <subgenre value="periodic">
                <genre-descr lang="ru" title="Журналы, газеты"/>
            </subgenre>
            <subgenre value="comics">
                <genre-descr lang="ru" title="Комиксы"/>
            </subgenre>
            <subgenre value="unfinished">
                <genre-descr lang="ru" title="Незавершенное"/>
            </subgenre>
            <subgenre value="other">
                <genre-descr lang="ru" title="Неотсортированное"/>
            </subgenre>
            <subgenre value="network_literature">
                <genre-descr lang="ru" title="Самиздат, сетевая литература"/>
            </subgenre>
            <subgenre value="fanfiction">
                <genre-descr lang="ru" title="Фанфик"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="religion">
        <root-descr lang="ru" genre-title="Религия, духовность, эзотерика" detailed=""/>
        <subgenres>
            <subgenre value="astrology">
                <genre-descr lang="ru" title="Астрология и хиромантия"/>
            </subgenre>
            <subgenre value="religion_budda">
                <genre-descr lang="ru" title="Буддизм"/>
            </subgenre>
            <subgenre value="religion_hinduism">
                <genre-descr lang="ru" title="Индуизм"/>
            </subgenre>
            <subgenre value="religion_islam">
                <genre-descr lang="ru" title="Ислам"/>
            </subgenre>
            <subgenre value="religion_judaism">
                <genre-descr lang="ru" title="Иудаизм"/>
            </subgenre>
            <subgenre value="religion_catholicism">
                <genre-descr lang="ru" title="Католицизм"/>
            </subgenre>
            <subgenre value="religion_orthodoxy">
                <genre-descr lang="ru" title="Православие"/>
            </subgenre>
            <subgenre value="religion_protestantism">
                <genre-descr lang="ru" title="Протестантизм"/>
            </subgenre>
            <subgenre value="sci_religion">
                <genre-descr lang="ru" title="Религиоведение"/>
            </subgenre>
            <subgenre value="religion">
                <genre-descr lang="ru" title="Религия, религиозная литература"/>
            </subgenre>
            <subgenre value="religion_self">
                <genre-descr lang="ru" title="Самосовершенствование"/>
            </subgenre>
            <subgenre value="religion_christianity">
                <genre-descr lang="ru" title="Христианство"/>
            </subgenre>
            <subgenre value="religion_esoterics">
                <genre-descr lang="ru" title="Эзотерика, эзотерическая литература"/>
            </subgenre>
            <subgenre value="religion_paganism">
                <genre-descr lang="ru" title="Язычество"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="reference">
        <root-descr lang="ru" genre-title="Справочная литература" detailed=""/>
        <subgenres>
            <subgenre value="geo_guides">
                <genre-descr lang="ru" title="Путеводители, карты, атласы"/>
            </subgenre>
            <subgenre value="ref_guide">
                <genre-descr lang="ru" title="Руководства"/>
            </subgenre>
            <subgenre value="ref_dict">
                <genre-descr lang="ru" title="Словари"/>
            </subgenre>
            <subgenre value="reference">
                <genre-descr lang="ru" title="Справочная литература"/>
            </subgenre>
            <subgenre value="ref_ref">
                <genre-descr lang="ru" title="Справочники"/>
            </subgenre>
            <subgenre value="ref_encyc">
                <genre-descr lang="ru" title="Энциклопедии"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="antique">
        <root-descr lang="ru" genre-title="Старинное" detailed=""/>
        <subgenres>
            <subgenre value="antique">
                <genre-descr lang="ru" title="Старинное"/>
            </subgenre>
            <subgenre value="antique_ant">
                <genre-descr lang="ru" title="Античная литератур"/>
            </subgenre>
            <subgenre value="antique_east">
                <genre-descr lang="ru" title="Древневосточная литература"/>
            </subgenre>
            <subgenre value="antique_russian">
                <genre-descr lang="ru" title="Древнерусская литература"/>
            </subgenre>
            <subgenre value="antique_european">
                <genre-descr lang="ru" title="Европейская старинная литература"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="sci_tech">
        <root-descr lang="ru" genre-title="Техника" detailed=""/>
        <subgenres>
            <subgenre value="auto_business">
                <genre-descr lang="ru" title="Автодело"/>
            </subgenre>
            <subgenre value="military_weapon">
                <genre-descr lang="ru" title="Военное дело, военная техника и вооружение"/>
            </subgenre>
            <subgenre value="equ_history">
                <genre-descr lang="ru" title="История техники"/>
            </subgenre>
            <subgenre value="sci_metal">
                <genre-descr lang="ru" title="Металлургия"/>
            </subgenre>
            <subgenre value="sci_radio">
                <genre-descr lang="ru" title="Радиоэлектроника"/>
            </subgenre>
            <subgenre value="sci_build">
                <genre-descr lang="ru" title="Строительство и сопромат"/>
            </subgenre>
            <subgenre value="sci_tech">
                <genre-descr lang="ru" title="Технические науки"/>
            </subgenre>
            <subgenre value="sci_transport">
                <genre-descr lang="ru" title="Транспорт и авиация"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="sci_textbook">
        <root-descr lang="ru" genre-title="Учебники и пособия" detailed=""/>
        <subgenres>
            <subgenre value="sci_textbook">
                <genre-descr lang="ru" title="Учебники и пособия"/>
            </subgenre>
            <subgenre value="tbg_higher">
                <genre-descr lang="ru" title="Учебники и пособия ВУЗов"/>
            </subgenre>
            <subgenre value="tbg_secondary">
                <genre-descr lang="ru" title="Учебники и пособия для среднего и специального образования"/>
            </subgenre>
            <subgenre value="tbg_school">
                <genre-descr lang="ru" title="Школьные учебники и пособия, рефераты, шпаргалки"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="sf">
        <root-descr lang="ru" genre-title="Фантастика" detailed=""/>
        <subgenres>
            <subgenre value="sf_history">
                <genre-descr lang="ru" title="Альтернативная история, попаданцы"/>
            </subgenre>
            <subgenre value="sf_action">
                <genre-descr lang="ru" title="Боевая фантастика"/>
            </subgenre>
            <subgenre value="sf_heroic">
                <genre-descr lang="ru" title="Героическая фантастика"/>
            </subgenre>
            <subgenre value="sf_fantasy_city">
                <genre-descr lang="ru" title="Городское фэнтези"/>
            </subgenre>
            <subgenre value="sf_detective">
                <genre-descr lang="ru" title="Детективная фантастика"/>
            </subgenre>
            <subgenre value="sf_cyberpunk">
                <genre-descr lang="ru" title="Киберпанк"/>
            </subgenre>
            <subgenre value="sf_space">
                <genre-descr lang="ru" title="Космическая фантастика"/>
            </subgenre>
            <subgenre value="sf_mystic">
                <genre-descr lang="ru" title="Мистика"/>
            </subgenre>
            <subgenre value="fairy_fantasy">
                <genre-descr lang="ru" title="Мифологическое фэнтези"/>
            </subgenre>
            <subgenre value="sf">
                <genre-descr lang="ru" title="Научная Фантастика"/>
            </subgenre>
            <subgenre value="sf_postapocalyptic">
                <genre-descr lang="ru" title="Постапокалипсис"/>
            </subgenre>
            <subgenre value="russian_fantasy">
                <genre-descr lang="ru" title="Славянское фэнтези"/>
            </subgenre>
            <subgenre value="modern_tale">
                <genre-descr lang="ru" title="Современная сказка"/>
            </subgenre>
            <subgenre value="sf_social">
                <genre-descr lang="ru" title="Социально-психологическая фантастика"/>
            </subgenre>
            <subgenre value="sf_stimpank">
                <genre-descr lang="ru" title="Стимпанк"/>
            </subgenre>
            <subgenre value="sf_technofantasy">
                <genre-descr lang="ru" title="Технофэнтези"/>
            </subgenre>
            <subgenre value="sf_horror">
                <genre-descr lang="ru" title="Ужасы"/>
            </subgenre>
            <subgenre value="sf_etc">
                <genre-descr lang="ru" title="Фантастика"/>
            </subgenre>
            <subgenre value="sf_fantasy">
                <genre-descr lang="ru" title="Фэнтези"/>
            </subgenre>
            <subgenre value="hronoopera">
                <genre-descr lang="ru" title="Хроноопера"/>
            </subgenre>
            <subgenre value="sf_epic">
                <genre-descr lang="ru" title="Эпическая фантастика"/>
            </subgenre>
            <subgenre value="theasf_humortre">
                <genre-descr lang="ru" title="Юмористическая фантастика"/>
            </subgenre>
            <subgenre value="theatre">
                <genre-descr lang="ru" title="Театр"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="folklore">
        <root-descr lang="ru" genre-title="Фольклор" detailed=""/>
        <subgenres>
            <subgenre value="epic">
                <genre-descr lang="ru" title="Былины, эпопея"/>
            </subgenre>
            <subgenre value="child_folklore">
                <genre-descr lang="ru" title="Детский фольклор"/>
            </subgenre>
            <subgenre value="antique_myths">
                <genre-descr lang="ru" title="Мифы. Легенды. Эпос"/>
            </subgenre>
            <subgenre value="folk_songs">
                <genre-descr lang="ru" title="Народные песни"/>
            </subgenre>
            <subgenre value="folk_tale">
                <genre-descr lang="ru" title="Народные сказки"/>
            </subgenre>
            <subgenre value="proverbs">
                <genre-descr lang="ru" title="Пословицы, поговорки"/>
            </subgenre>
            <subgenre value="folklore">
                <genre-descr lang="ru" title="Фольклор, загадки folklore"/>
            </subgenre>
            <subgenre value="limerick">
                <genre-descr lang="ru" title="Частушки, прибаутки, потешки"/>
            </subgenre>
        </subgenres>
    </genre>
    <genre value="humor">
        <root-descr lang="ru" genre-title="Юмор" detailed=""/>
        <subgenres>
            <subgenre value="humor_anecdote">
                <genre-descr lang="ru" title="Анекдоты"/>
            </subgenre>
            <subgenre value="humor_satire">
                <genre-descr lang="ru" title="Сатира"/>
            </subgenre>
            <subgenre value="humor">
                <genre-descr lang="ru" title="Юмор"/>
            </subgenre>
            <subgenre value="humor_prose">
                <genre-descr lang="ru" title="Юмористическая проза"/>
            </subgenre>
        </subgenres>
    </genre>
</fbgenrestransfer>
`
