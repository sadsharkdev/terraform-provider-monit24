# Monit24 API v3.51 - Full Endpoint Inventory

Generated 2026-08-08 from `https://api.monit24.pl/v3/public/swagger.json` (OpenAPI 3, 361 paths / 441 operations, 35 tags). Reference material for planning provider updates — not part of the Terraform-facing `docs/`. See the provider's `CLAUDE.md` "Resources" section for what's actually implemented.

## accounts

- `GET /accounts` — Zwraca listę kont.
- `POST /accounts` — Tworzy dodatkowego użytkownika konta.
- `POST /accounts/register` — Zakłada nowe konto.
- `POST /accounts/search` — Zaawansowane pobieranie kont.
- `POST /accounts/subaccount` — Tworzy nowe konto zależne, z możliwością udostępniania grup.
- `DELETE /accounts/{_id}` — Usuwa dodatkowego użytkownika konta.
- `GET /accounts/{_id}` — Zwraca pojedyncze konto.
- `PUT /accounts/{_id}` — Modyfikuje ustawienia konta.
- `POST /accounts/{_id}/2fa/setup` — Generuje i zwraca klucz 2FA dla konta. Aby zakończyć konfigurację 2FA należy na jego podstawie wygenerować token TOTP i potwierdzić operacją 2fa/verify.
- `POST /accounts/{_id}/2fa/verify` — Weryfikuje jednorazowy kod 2FA i włącza 2FA dla konta.
- `GET /accounts/{_id}/activate` — Aktywuje konto.
- `POST /accounts/{_id}/block` — Blokuje konto.
- `POST /accounts/{_id}/change_password` — Zmienia hasło do konta.
- `POST /accounts/{_id}/close` — Zamyka konto.
- `POST /accounts/{_id}/copy_group_shares` — Kopiuje uprawnienia do grup z jednego konta do wielu.
- `POST /accounts/{_id}/customer_portal_sessions` — Tworzy nową sesję portalu billingowego.
- `GET /accounts/{_id}/language` — Zwraca domyślny język powiadomień i raportów dla konta.
- `GET /accounts/{_id}/stats` — Zwraca statystyki konta.
- `GET /accounts/{_id}/stats/by_service_type_category` — Zwraca statystyki konta z podziałem na kategorie typów usług.
- `GET /accounts/{_id}/time_zone` — Zwraca domyślną strefę czasową zawieszeń tygodniowych i trybu cichego dla konta.
- `POST /accounts/{_id}/unblock` — Odblokowuje konto.
- `GET /contact_groups/{_id}/owner` — Zwraca konto, do którego należy wybrana grupa kontaktów.
- `GET /contacts/{_id}/owner` — Zwraca konto, do którego należy wybrany kontakt.
- `GET /events/{_id}/owner` — Zwraca konto, do którego należy wybrane zdarzenie.
- `GET /groups/{_id}/owner` — Zwraca konto, do którego należy wybrana grupa.
- `GET /groups/{_id}/shares` — Zwraca ustawienia udostępniania grupy.
- `DELETE /groups/{_id}/shares/{_account_id}` — Usuwa udostępnienie grupy użytkownikowi.
- `GET /groups/{_id}/shares/{_account_id}` — Zwraca ustawienia udostępnienia grupy wybranemu użytkownikowi.
- `PUT /groups/{_id}/shares/{_account_id}` — Tworzy lub modyfikuje udostępnienie grupy użytkownikowi.
- `GET /groups/{_id}/shares/{_account_id}/permissions` — Zwraca uprawnienia do ustawień udostępnienia grupy wybranemu użytkownikowi.
- `GET /notification_addresses/{_id}/owner` — Zwraca konto, do którego należy adres powiadomień.
- `GET /periodic_report_addresses/{_id}/owner` — Zwraca konto, do którego należy adres raportów okresowych.
- `GET /report_templates/{_id}/owner` — Zwraca konto, do którego należy szablon.
- `GET /reports/{_id}/owner` — Zwraca konto, do którego należy raport.
- `GET /services/{_id}/owner` — Zwraca konto, do którego należy wybrana usługa.
- `GET /templates/{_id}/owner` — Zwraca konto, do którego należy wybrany szablon.

## admin

- `GET /admin/accounts` — Zwraca listę kont.
- `POST /admin/accounts` — Tworzy nowe konto.
- `POST /admin/accounts/search` — Zaawansowane pobieranie kont.
- `DELETE /admin/accounts/{_id}` — Usuwa konto.
- `GET /admin/accounts/{_id}` — Zwraca pojedyncze konto.
- `PUT /admin/accounts/{_id}` — Modyfikuje konto.
- `POST /admin/accounts/{_id}/block` — Blokuje konto.
- `POST /admin/accounts/{_id}/unblock` — Odblokowuje konto.
- `GET /admin/packages` — Zwraca listę pakietów.
- `POST /admin/packages/search` — Zaawansowane pobieranie pakietów.
- `GET /admin/packages/{_id}` — Zwraca pojedynczy pakiet.
- `GET /admin/usage_stats/by_account/{year}/{month}` — Statystyki zużycia zasobów dla kont.
- `GET /admin/usage_stats/by_group/{year}/{month}` — Statystyki zużycia zasobów dla grup.

## auth_token

- `DELETE /auth_token` **[DEPRECATED]** — Kończy sesję.
- `GET /auth_token` **[DEPRECATED]** — Zwraca informacje o sesji.
- `POST /auth_token` **[DEPRECATED]** — Rozpoczyna nową sesję.

## charts

- `GET /services/{_id}/charts/availability` — Zwraca dane do wykresu Service Level.
- `GET /services/{_id}/charts/availability/range` — Zwraca przedział czasowy, z którego są dostępne dane o dostępności usługi.
- `GET /services/{_id}/charts/performance` — Zwraca dane do wykresu czasu odpowiedzi, z automatyczną agregacją.
- `GET /services/{_id}/charts/performance/checks` — Zwraca dane do wykresu czasu odpowiedzi pojedynczych sprawdzeń.
- `GET /services/{_id}/charts/performance/range` — Zwraca przedział czasowy, z którego są dostępne dane do wykresu czasu odpowiedzi.
- `GET /services/{_id}/charts/performance/river` — Zwraca zagregowane dane do wykresu czasu odpowiedzi.
- `GET /services/{_id}/charts/performance/steps` — Zwraca dane do wykresu czasu odpowiedzi, z podziałem na kroki.
- `GET /services/{_id}/charts/performance/used_sensors` — Zwraca stacje monitorujące, z których usługa była monitorowana.
- `GET /services/{_id}/event_charts/{_event_id}/availability` — Zwraca dane do wykresu Service Level dla wybranego zdarzenia.
- `GET /services/{_id}/event_charts/{_event_id}/availability/range` — Zwraca przedział czasowy, z którego są dostępne dane Service Level dla zdarzenia.

## contact_addresses

- `GET /contact_addresses` — Zwraca listę adresów.
- `POST /contact_addresses` — Tworzy nowy adres.
- `POST /contact_addresses/search` — Zaawansowane pobieranie adresów.
- `DELETE /contact_addresses/{_id}` — Usuwa adres.
- `GET /contact_addresses/{_id}` — Zwraca pojedynczy adres.
- `PUT /contact_addresses/{_id}` — Modyfikuje adres.
- `GET /contact_addresses/{_id}/channel` — Zwraca kanał powiadomień dla adresu.
- `GET /contacts/{_id}/addresses` — Zwraca listę adresów należących do kontaktu.

## contact_groups

- `GET /contact_groups` — Zwraca listę grup kontaktów.
- `POST /contact_groups` — Tworzy nową grupę kontaktów.
- `POST /contact_groups/search` — Zaawansowane pobieranie grup kontaktów.
- `DELETE /contact_groups/{_id}` — Usuwa grupę kontaktów.
- `GET /contact_groups/{_id}` — Zwraca pojedynczą grupę kontaktów.
- `PUT /contact_groups/{_id}` — Modyfikuje grupę kontaktów.

## contact_suspensions

- `GET /contact_suspensions` — Zwraca listę zawieszeń.
- `POST /contact_suspensions` — Tworzy nowe zawieszenie.
- `POST /contact_suspensions/search` — Zaawansowane pobieranie zawieszeń.
- `DELETE /contact_suspensions/{_id}` — Usuwa zawieszenie.
- `GET /contact_suspensions/{_id}` — Zwraca pojedyncze zawieszenie.
- `PUT /contact_suspensions/{_id}` — Modyfikuje zawieszenie.
- `GET /contact_weekly_suspensions` — Zwraca listę tygodniowych zawieszeń konktaktów.
- `POST /contact_weekly_suspensions` — Tworzy nowe tygodniowe zawieszenie kontaktu.
- `POST /contact_weekly_suspensions/search` — Zaawansowane pobieranie tygodniowych zawieszeń kontaktów.
- `DELETE /contact_weekly_suspensions/{_id}` — Usuwa tygodniowe zawieszenie kontaktu.
- `GET /contact_weekly_suspensions/{_id}` — Zwraca pojedyncze tygodniowe zawieszenie kontaktu.
- `PUT /contact_weekly_suspensions/{_id}` — Modyfikuje tygodniowe zawieszenie kontaktu.
- `GET /contact_weekly_suspensions/{_id}/permissions` — Zwraca uprawnienia do wybranego zawieszenia tygodniowego.
- `GET /contacts/{_id}/suspensions` — Zwraca listę zawieszeń kontaktu.
- `GET /contacts/{_id}/weekly_suspensions` — Zwraca listę tygodniowych zawieszeń kontaktu.

## contacts

- `GET /accounts/{_id}/contact_groups` — Zwraca listę grup kontaktów należących do konta.
- `GET /accounts/{_id}/contacts` — Zwraca listę kontaktów należących do konta.
- `GET /contact_addresses/{_id}/contact` — Zwraca kontakt, do którego należy adres.
- `GET /contact_groups/{_id}/contacts` — Zwraca listę kontaktów należących do grupy.
- `GET /contact_suspensions/{_id}/contact` — Zwraca kontakt, do którego należy zawieszenie.
- `GET /contact_weekly_suspensions/{_id}/contact` — Zwraca kontakt, którego dotyczy wybrane tygodniowe zawieszenie.
- `GET /contacts` — Zwraca listę kontaktów.
- `POST /contacts` — Tworzy nowy kontakt.
- `POST /contacts/search` — Zaawansowane pobieranie kontaktów.
- `DELETE /contacts/{_id}` — Usuwa kontakt.
- `GET /contacts/{_id}` — Zwraca pojedynczy kontakt.
- `PUT /contacts/{_id}` — Modyfikuje kontakt.

## corrections

- `GET /corrections` — Zwraca listę korekt.
- `POST /corrections` — Tworzy nową korektę.
- `POST /corrections/bulk_create` — Dodaje korektę do wielu usług jednocześnie.
- `POST /corrections/bulk_delete` — Usuwa wiele korekt jednocześnie.
- `POST /corrections/bulk_update` — Modyfikuje wiele korekt jednocześnie.
- `POST /corrections/search` — Zaawansowane pobieranie korekt.
- `DELETE /corrections/{_id}` — Usuwa korektę.
- `GET /corrections/{_id}` — Zwraca pojedynczą korektę.
- `PUT /corrections/{_id}` — Modyfikuje istniejącą korektę.
- `GET /services/{_id}/corrections` — Zwraca listę korekt Service Level usługi.
- `POST /services/{_id}/corrections/bulk_create` — Tworzy wiele korekt SL usługi jednocześnie.

## dictionaries

- `GET /dictionaries/artifact_types` — Zwraca listę znanych typów artefaktów.
- `GET /dictionaries/artifact_types/{_id}` — Zwraca wybrany typ artefaktu.
- `GET /dictionaries/contact_channels` — Zwraca listę kanałów kontaktowych.
- `GET /dictionaries/contact_channels/{_id}` — Zwraca wybrany kanał kontaktowy.
- `GET /dictionaries/error_type_categories` — Zwraca listę kategorii typów błędu.
- `GET /dictionaries/error_type_categories/{_id}` — Zwraca wybraną kategorię typów błędu.
- `GET /dictionaries/languages` — Zwraca listę dostępnych języków.
- `GET /dictionaries/languages/{_id}` — Zwraca wybrany język.
- `GET /dictionaries/metric_types` — Zwraca listę znanych typów metryk.
- `GET /dictionaries/metric_types/{_id}` — Zwraca wybrany typ metryki.
- `GET /dictionaries/notification_channels` — Zwraca listę kanałów powiadomień dostępnych w systemie.
- `GET /dictionaries/notification_channels/{_id}` — Zwraca wybrany kanał powiadomień.
- `GET /dictionaries/notification_conditions` — Zwraca listę dostępnych warunków powiadomień.
- `GET /dictionaries/notification_conditions/{_id}` — Zwraca wybrany warunek powiadomień.
- `GET /dictionaries/notification_modes` — Zwraca listę trybów powiadomień.
- `GET /dictionaries/notification_modes/{_id}` — Zwraca wybrany tryb powiadomień.
- `GET /dictionaries/recovery_notification_modes` — Zwraca listę trybów powiadomień o końcu awarii.
- `GET /dictionaries/recovery_notification_modes/{_id}` — Zwraca wybrany tryb powiadomień o końcu awarii.
- `GET /dictionaries/report_attachment_types` — Zwraca listę znanych typów załączników do raportów.
- `GET /dictionaries/report_attachment_types/{_id}` — Zwraca wybrany typ załcznika do raportu.
- `GET /dictionaries/report_sections` — Zwraca listę dostępnych sekcji raportów.
- `GET /dictionaries/report_sections/{_id}` — Zwraca wybraną sekcję raportów.
- `GET /dictionaries/report_template_variables` — Zwraca listę zmiennych dostępnych w szablonach raportów.
- `GET /dictionaries/report_template_variables/{_id}` — Zwraca wybraną zmienną dostępną w szablonach raportów.
- `GET /dictionaries/reporting_plans` — Zwraca listę dostępnych planów generowania raportów.
- `GET /dictionaries/reporting_plans/{_id}` — Zwraca wybrany plan generowania raportów.
- `GET /dictionaries/service_type_categories` — Zwraca listę kategorii typów usług.
- `GET /dictionaries/service_type_categories/{_id}` — Zwraca wybraną kategorię typów usług.
- `GET /dictionaries/template_variables` — Zwraca listę zmiennych dostępnych w szablonach powiadomień.
- `GET /dictionaries/template_variables/{_id}` — Zwraca wybraną zmienną dostępną w szablonach powiadomień.
- `GET /dictionaries/time_zones` — Zwraca listę dostępnych stref czasowych.
- `GET /dictionaries/time_zones/{_id}` — Zwraca wybraną strefę czasową.

## error_types

- `GET /dictionaries/error_type_categories/{_id}/error_types` — Zwraca listę typów błędów należących do wybranej kategorii.
- `GET /error_types` — Zwraca listę typów błędów.
- `POST /error_types/search` — Zaawansowane pobieranie typów błędów.
- `GET /error_types/{_id}` — Zwraca pojedynczy typ błędu.

## escalation_suspensions

- `GET /escalation_suspensions` — Zwraca listę zawieszeń.
- `POST /escalation_suspensions` — Tworzy nowe zawieszenie.
- `POST /escalation_suspensions/search` — Zaawansowane pobieranie zawieszeń.
- `DELETE /escalation_suspensions/{_id}` — Usuwa zawieszenie.
- `GET /escalation_suspensions/{_id}` — Zwraca pojedyncze zawieszenie.
- `PUT /escalation_suspensions/{_id}` — Modyfikuje zawieszenie.
- `GET /escalation_weekly_suspensions` — Zwraca listę tygodniowych zawieszeń eskalacji.
- `POST /escalation_weekly_suspensions` — Tworzy nowe tygodniowe zawieszenie eskalacji.
- `POST /escalation_weekly_suspensions/search` — Zaawansowane pobieranie tygodniowych zawieszeń eskalacji.
- `DELETE /escalation_weekly_suspensions/{_id}` — Usuwa tygodniowe zawieszenie eskalacji.
- `GET /escalation_weekly_suspensions/{_id}` — Zwraca pojedyncze tygodniowe zawieszenie eskalacji.
- `PUT /escalation_weekly_suspensions/{_id}` — Modyfikuje tygodniowe zawieszenie eskalacji.
- `GET /escalation_weekly_suspensions/{_id}/permissions` — Zwraca uprawnienia do wybranego zawieszenia tygodniowego.
- `GET /escalations/{_id}/suspensions` — Zwraca listę zawieszeń eskalacji.
- `GET /escalations/{_id}/weekly_suspensions` — Zwraca listę tygodniowych zawieszeń eskalacji.

## escalations

- `GET /escalation_statuses` — Zwraca listę aktualnych statusów eskalacji.
- `POST /escalation_statuses/search` — Zaawansowane pobieranie statusów eskalacji.
- `GET /escalation_suspensions/{_id}/escalation` — Zwraca eskalację, której dotyczy zawieszenie.
- `GET /escalation_weekly_suspensions/{_id}/escalation` — Zwraca eskalację, której dotyczy wybrane tygodniowe zawieszenie.
- `GET /escalations` — Zwraca listę eskalacji.
- `POST /escalations` — Tworzy nową eskalację.
- `POST /escalations/search` — Zaawansowane pobieranie eskalacji.
- `DELETE /escalations/{_id}` — Usuwa eskalację.
- `GET /escalations/{_id}` — Zwraca pojedynczą eskalację.
- `PUT /escalations/{_id}` — Modyfikuje eskalację.
- `GET /escalations/{_id}/statuses` — Zwraca listę aktualnych statusów eskalacji dla usług.
- `GET /events/{_id}/escalations` — Zwraca listę eskalacji przypisanych do zdarzenia.
- `GET /services/{_id}/escalation_statuses/{_escalation_id}` — Zwraca aktualny status pojedynczej eskalacji dla usługi.

## events

- `GET /accounts/{_id}/events` — Zwraca listę zdarzeń należących do konta.
- `GET /escalations/{_id}/event` — Zwraca zdarzenie, o którym należy powiadamiać.
- `GET /event_statuses` — Zwraca listę aktualnych statusów zdarzeń.
- `POST /event_statuses/search` — Zaawansowane pobieranie statusów zdarzeń.
- `GET /events` — Zwraca listę zdarzeń.
- `POST /events` — Tworzy nowe zdarzenie.
- `POST /events/search` — Zaawansowane pobieranie zdarzeń.
- `DELETE /events/{_id}` — Usuwa zdarzenie razem z jego eskalacjami.
- `GET /events/{_id}` — Zwraca pojedyncze zdarzenie.
- `PUT /events/{_id}` — Modyfikuje zdarzenie.
- `GET /events/{_id}/possible_conditions` — Zwraca listę zdarzeń od których może zależeć wybrane zdarzenie.
- `GET /events/{_id}/statuses` — Zwraca listę aktualnych statusów zdarzenia dla poszczególnych usług.
- `GET /services/{_id}/event_statuses/{_event_id}` — Zwraca aktualny status pojedynczego zdarzenia dla usługi.
- `POST /services/{_id}/event_statuses/{_event_id}/disable` — Wyłącza zachodzenie zdarzenia.

## group_stats

- `GET /group_stats` — Zwraca listę statystyk grup.
- `POST /group_stats/search` — Zaawansowane pobieranie statystyk grup.
- `GET /group_stats/{_id}` — Zwraca aktualne statystyki pojedynczej grupy.

## groups

- `GET /accounts/{_id}/default_group` — Zwraca domyślną grupę wybranego konta.
- `GET /accounts/{_id}/group_shares` — Zwraca listę udostępnień grup należących do konta.
- `GET /accounts/{_id}/groups` — Zwraca listę grup należących do konta.
- `GET /groups` — Zwraca listę grup usług.
- `POST /groups` — Tworzy nową grupę.
- `POST /groups/bulk_delete` — Usuwa wiele grup jednocześnie.
- `POST /groups/search` — Zaawansowane pobieranie grup usług.
- `DELETE /groups/{_id}` — Usuwa grupę usług.
- `GET /groups/{_id}` — Zwraca pojedynczą grupę usług.
- `PUT /groups/{_id}` — Modyfikuje grupę.
- `POST /groups/{_id}/send_test_sms_notification` — Wysyła testowe powiadomienie SMS na numery telefonu przypisane do grupy.
- `GET /notification_addresses/{_id}/group` — Zwraca grupę, do której jest przypisany adres powiadomień.
- `GET /periodic_report_addresses/{_id}/group` — Zwraca grupę, do której jest przypisany adres raportów okresowych.
- `GET /services/{_id}/available_groups` — Zwraca listę grup, do których można przenieść wybraną usługę.
- `GET /services/{_id}/group` — Zwraca grupę, do której należy wybrana usługa.

## history

- `GET /services/{_id}/event_history/{_event_id}/service_level` — Zwraca Service Level dla danej usługi i zdarzenia w wybranym przedziale czasowym.
- `GET /services/{_id}/event_history/{_event_id}/status_changes` — Zwraca historię zmian stanów zdarzenia dla usługi z wybranego przedziału czasowego.
- `GET /services/{_id}/event_history/{_event_id}/status_changes/{_change_id}` — Zwraca wybraną zmianę stanu zdarzenia dla usługi.
- `GET /services/{_id}/history/analyses` — Zwraca historię analiz działania usługi z wybranego przedziału czasowego.
- `GET /services/{_id}/history/analyses/{_analysis_id}` — Zwraca pojedynczą analizę działania usługi.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks` — Zwraca listę sprawdzeń usługi należących do wybranej analizy.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}` — Zwraca pojedyncze sprawdzenie usługi.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/artifacts` — Zwraca listę artefaktów dla wybranego sprawdzenia.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/artifacts/{_artifact_id}` — Zwraca informacje o wybranym artefakcie.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/artifacts/{_artifact_id}/by_token/{file_name}` — Zwraca zawartość wybranego artefaktu na podstawie tokenu.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/artifacts/{_artifact_id}/plain` — Zwraca zawartość wybranego artefaktu.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/artifacts/{_artifact_id}/zip` — Zwraca zawartość wybranego artefaktu w postaci archiwum ZIP.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/content` **[DEPRECATED]** — Zwraca treść strony z momentu awarii.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/details` — Zwraca szczegóły pojedynczego sprawdzenia usługi.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/har` **[DEPRECATED]** — Zwraca plik HAR dla sprawdzenia.
- `GET /services/{_id}/history/analyses/{_analysis_id}/checks/{_sensor_id}/uncompressed_content` **[DEPRECATED]** — Zwraca treść strony z momentu awarii (bez kompresji).
- `GET /services/{_id}/history/analyses/{_analysis_id}/steps_summary` — Zwraca statystyki kroków scenariusza dla wybranej analizy.
- `GET /services/{_id}/history/checks` — Zwraca historię sprawdzeń usługi z wybranego przedziału czasowego.
- `GET /services/{_id}/history/detected_error_types` — Zwraca błędy wykryte podczas monitoringu usługi.
- `GET /services/{_id}/history/service_level` — Zwraca dostępność usługi w wybranym przedziale czasowym.
- `GET /services/{_id}/history/status_changes` — Zwraca historię zmian stanów działania usługi z wybranego przedziału czasowego.
- `POST /services/{_id}/history/status_changes/bulk_update` — Aktualizuje opisy wielu zmian stanów jednocześnie.
- `GET /services/{_id}/history/status_changes/{_change_id}` — Zwraca wybraną zmianę stanu działania usługi.
- `PUT /services/{_id}/history/status_changes/{_change_id}` — Modyfikuje wybraną zmianę stanu działania usługi.
- `GET /services/{_id}/history/status_changes/{_change_id}/permissions` — Zwraca uprawnienia do historycznej zmiany stanu działania usługi.
- `GET /services/{_id}/history/status_changes/{_change_id}/stats` — Zwraca statystyki liczby sprawdzeń w wybranym stanie działania usługi.
- `GET /services/{_id}/history/time_range_info` — Zwraca informacje o dostępności historii w danym przedziale czasowym.
- `GET /services/{_id}/history/used_sensors` — Zwraca stacje monitorujące, z których usługa była monitorowana.

## logs

- `GET /accounts/{_account_id}/logs/authentications` — Zwraca logi pomyślnych autentykacji do konta.
- `GET /accounts/{_account_id}/logs/contact_notifications` — Zwraca logi powiadomień wysłanych do kontaktów należących do konta.
- `GET /accounts/{_account_id}/logs/contact_notifications/{_id}` — Zwraca pojedyncze powiadomienie wysłane do kontaktu.
- `GET /accounts/{_account_id}/logs/event_status_changes` — Zwraca logi zmian stanów zdarzeń należących do konta.
- `GET /accounts/{_account_id}/logs/notifications` — Zwraca logi wysłanych powiadomień dla konta.
- `GET /accounts/{_account_id}/logs/notifications/{_id}` — Zwraca pojedyncze powiadomienie.
- `GET /accounts/{_account_id}/logs/status_changes` — Zwraca logi zmian stanów usług należących do konta.

## mobile

- `GET /groups/{_id}/mobile_devices_count` — Zwraca liczbę urządzeń mobilnych otrzymujących powiadomienia o usługach z grupy.
- `DELETE /mobile/devices/{_id}` — Wyrejestrowuje urządzenie mobilne z powiadomień push.
- `PUT /mobile/devices/{_id}` — Rejestruje urządzenie mobilne do powiadomień push.
- `GET /mobile/notifications/{_id}` — Zwraca pełną treść powiadomienia push.
- `GET /services/{_id}/mobile_devices_count` — Zwraca liczbę urządzeń mobilnych otrzymujących powiadomienia o usłudze.

## notification_addresses

- `GET /accounts/{_id}/notification_addresses` — Zwraca listę adresów powiadomień należących do konta.
- `GET /groups/{_id}/notification_addresses` — Zwraca listę adresów powiadomień przypisanych do grupy.
- `GET /notification_addresses` — Zwraca listę adresów powiadomień.
- `POST /notification_addresses` — Tworzy nowy adres powiadomień.
- `POST /notification_addresses/search` — Zaawansowane pobieranie adresów powiadomień.
- `DELETE /notification_addresses/{_id}` — Usuwa adres powiadomień.
- `GET /notification_addresses/{_id}` — Zwraca pojedynczy adres powiadomień.
- `PUT /notification_addresses/{_id}` — Modyfikuje adres powiadomień.
- `GET /notification_addresses/{_id}/notification_channel` — Zwraca kanał powiadomień dla wybranego adresu.
- `POST /notification_addresses/{_id}/send_test_sms_notification` — Wysyła testowe powiadomienie SMS.
- `GET /services/{_id}/notification_addresses` — Zwraca listę adresów powiadomień włączonych dla usługi.

## packages

- `GET /accounts/{_id}/package` — Zwraca pakiet przypisany do konta.
- `GET /packages` — Zwraca listę pakietów.
- `GET /packages/available_for_new_accounts` — Zwraca listę pakietów dostępnych przy rejestracji nowego konta.
- `POST /packages/search` — Zaawansowane pobieranie pakietów.
- `GET /packages/{_id}` — Zwraca pojedynczy pakiet.
- `GET /packages/{_id}/available_notification_channels` — Zwraca listę kanałów powiadomień dostępnych w pakiecie.

## periodic_report_addresses

- `GET /accounts/{_id}/periodic_report_addresses` — Zwraca listę adresów raportów okresowych należących do konta.
- `GET /groups/{_id}/periodic_report_addresses` — Zwraca listę adresów raportów okresowych przypisanych do grupy.
- `GET /periodic_report_addresses` — Zwraca listę adresów raportów okresowych.
- `POST /periodic_report_addresses` — Tworzy nowy adres raportów okresowych.
- `POST /periodic_report_addresses/search` — Zaawansowane pobieranie adresów raportów okresowych.
- `DELETE /periodic_report_addresses/{_id}` — Usuwa adres raportów okresowych.
- `GET /periodic_report_addresses/{_id}` — Zwraca pojedynczy adres raportów okresowych.
- `PUT /periodic_report_addresses/{_id}` — Modyfikuje adres raportów okresowych.
- `GET /services/{_id}/periodic_report_addresses` — Zwraca listę adresów, na które są wysyłane raporty okresowe dla usługi.

## permissions

- `GET /accounts/{_id}/permissions` — Zwraca uprawnienia do wybranego konta.
- `GET /admin/accounts/{_id}/permissions` — Zwraca uprawnienia do wybranego konta.
- `GET /admin/packages/{_id}/permissions` — Zwraca uprawnienia do wybranego pakietu.
- `GET /contact_addresses/{_id}/permissions` — Zwraca uprawnienia do wybranego adresu.
- `GET /contact_groups/{_id}/permissions` — Zwraca uprawnienia do wybranej grupy kontaktów.
- `GET /contact_suspensions/{_id}/permissions` — Zwraca uprawnienia do wybranego zawieszenia.
- `GET /contacts/{_id}/permissions` — Zwraca uprawnienia do wybranego kontaktu.
- `GET /corrections/{_id}/permissions` — Zwraca uprawnienia do wybranej korekty.
- `GET /error_types/{_id}/permissions` — Zwraca uprawnienia do wybranego typu błędu.
- `GET /escalation_suspensions/{_id}/permissions` — Zwraca uprawnienia do wybranego zawieszenia.
- `GET /escalations/{_id}/permissions` — Zwraca uprawnienia do wybranej eskalacji.
- `GET /events/{_id}/permissions` — Zwraca uprawnienia do wybranego zdarzenia.
- `GET /group_stats/{_id}/permissions` — Zwraca uprawnienia do statystyk wybranej grupy.
- `GET /groups/{_id}/permissions` — Zwraca uprawnienia do wybranej grupy usług.
- `GET /notification_addresses/{_id}/permissions` — Zwraca uprawnienia do adresu powiadomień.
- `GET /packages/{_id}/permissions` — Zwraca uprawnienia do wybranego pakietu.
- `GET /periodic_report_addresses/{_id}/permissions` — Zwraca uprawnienia do adresu raportów okresowych.
- `GET /report_templates/{_id}/permissions` — Zwraca uprawnienia do szablonu.
- `GET /reports/{_id}/permissions` — Zwraca uprawnienia do raportu.
- `GET /sensors/{_id}/permissions` — Zwraca uprawnienia do wybranej stacji monitorującej.
- `GET /service_statuses/{_id}/permissions` — Zwraca uprawnienia do aktualnego statusu wybranej usługi.
- `GET /service_types/{_id}/permissions` — Zwraca uprawnienia do wybranego typu usługi.
- `GET /services/{_id}/permissions` — Zwraca uprawnienia do wybranej usługi.
- `GET /sessions/{_id}/permissions` — Zwraca uprawnienia do wybranej sesji API.
- `GET /suspensions/{_id}/permissions` — Zwraca uprawnienia do wybranego zawieszenia.
- `GET /templates/{_id}/permissions` — Zwraca uprawnienia do wybranego szablonu.
- `GET /user_data/{_id}/permissions` — Zwraca uprawnienia do danych wybranego użytkownika.

## report_templates

- `GET /accounts/{_id}/report_templates` — Zwraca listę szablonów raportów należących do konta.
- `GET /report_templates` — Zwraca listę szablonów.
- `POST /report_templates` — Dodaje nowy szablon.
- `POST /report_templates/search` — Zaawansowane pobieranie szablonów.
- `DELETE /report_templates/{_id}` — Usuwa szablon razem z jego raportami.
- `GET /report_templates/{_id}` — Zwraca pojedynczy szablon.
- `PUT /report_templates/{_id}` — Modyfikuje szablon.
- `GET /report_templates/{_id}/reporting_plans` — Zwraca listę aktywnych planów generowania raportów według szablonu.
- `GET /report_templates/{_id}/sections` — Zwraca listę sekcji w raportach generowanych według szablonu.
- `GET /report_templates/{_id}/stats` — Zwraca statystyki wybranego szablonu.
- `GET /reports/{_id}/template` — Zwraca szablon, według którego raport został utworzony.

## reports

- `GET /accounts/{_id}/reports` — Zwraca listę raportów należących do konta.
- `GET /report_templates/{_id}/reports` — Zwraca listę raportów utworzonych według szablonu.
- `GET /reports` — Zwraca listę raportów.
- `POST /reports` — Zleca wygenerowanie nowego raportu na podstawie szablonu.
- `POST /reports/search` — Zaawansowane pobieranie raportów.
- `DELETE /reports/{_id}` — Usuwa raport.
- `GET /reports/{_id}` — Zwraca pojedynczy raport.
- `GET /reports/{_id}/generated_attachments` — Zwraca listę załączników do raportu.
- `GET /reports/{_id}/generated_attachments/{_type_id}/by_token/{file_name}` — Zwraca wybrany załącznik do raportu za pomocą tokena.
- `GET /reports/{_id}/generated_attachments/{_type_id}/plain` — Zwraca zawartość wybranego załącznika do raportu.
- `GET /reports/{_id}/json` — Zwraca plik z raportem w formacie JSON.
- `GET /reports/{_id}/pdf` — Zwraca plik z raportem w formacie PDF.
- `GET /reports/{_id}/pdf/by_token/{file_name}` — Zwraca plik z raportem w formacie PDF za pomocą tokena.
- `GET /reports/{_id}/sections` — Zwraca listę sekcji w raporcie.
- `GET /reports/{_report_id}/generated_attachments/{_type_id}` — Zwraca dane pojedynczego załącznika do raportu.

## sensors

- `GET /accounts/{_id}/available_sensors` — Zwraca listę stacji monitorujących dostępnych dla konta.
- `GET /dictionaries/service_type_categories/{_id}/sensors` — Zwraca listę stacji monitorujących obsługujących sprawdzanie usług wybranej kategorii.
- `GET /groups/{_id}/available_sensors` — Zwraca listę stacji monitorujących dostępnych dla grupy.
- `GET /groups/{_id}/available_sensors/by_service_type/{_service_type_id}` — Zwraca listę stacji monitorujących dostępnych dla usług danego typu w grupie.
- `GET /groups/{_id}/sensors` — Zwraca listę stacji monitorujących przypisanych do grupy w kategorii `default`.
- `GET /sensors` — Zwraca listę stacji monitorujących.
- `POST /sensors/search` — Zaawansowane pobieranie stacji monitorujących.
- `GET /sensors/{_id}` — Zwraca pojedynczą stację monitorującą.
- `GET /sensors/{_id}/status` — Zwraca aktualny status działania wybranej stacji monitorującej.
- `GET /services/{_id}/available_sensors` — Zwraca listę stacji monitorujących dostępnych dla usługi.
- `GET /services/{_id}/monitoring_sensors` — Zwraca listę stacji monitorujących, z których usługa jest monitorowana.

## service_statuses

- `GET /service_statuses` — Zwraca listę aktualnych statusów usług.
- `POST /service_statuses/search` — Zaawansowane pobieranie aktualnych statusów usług.
- `GET /service_statuses/{_id}` — Zwraca aktualny status pojedynczej usługi.

## service_suspensions

- `GET /scheduled_suspensions` **[DEPRECATED]** — Zwraca listę zaplanowanych zawieszeń.
- `POST /scheduled_suspensions` **[DEPRECATED]** — Tworzy nowe zawieszenie.
- `POST /scheduled_suspensions/bulk_create` **[DEPRECATED]** — Dodaje zaplanowane zawieszenie do wielu usług jednocześnie.
- `POST /scheduled_suspensions/bulk_delete` **[DEPRECATED]** — Usuwa wiele zawieszeń jednocześnie.
- `POST /scheduled_suspensions/search` **[DEPRECATED]** — Zaawansowane pobieranie zaplanowanych zawieszeń.
- `DELETE /scheduled_suspensions/{_id}` **[DEPRECATED]** — Usuwa zaplanowane zawieszenie.
- `GET /scheduled_suspensions/{_id}` **[DEPRECATED]** — Zwraca pojedyncze zaplanowane zawieszenie.
- `PUT /scheduled_suspensions/{_id}` **[DEPRECATED]** — Modyfikuje zaplanowane zawieszenie.
- `GET /scheduled_suspensions/{_id}/permissions` **[DEPRECATED]** — Zwraca uprawnienia do wybranego zawieszenia.
- `GET /services/{_id}/scheduled_suspensions` **[DEPRECATED]** — Zwraca listę zaplanowanych zawieszeń dla usługi.

## service_types

- `GET /accounts/{_id}/available_service_types` — Zwraca listę typów usług do wykorzystania w pakiecie dla konta.
- `GET /groups/{_id}/available_service_types` — Zwraca listę typów usług do wykorzystania dla konta, do którego należy grupa.
- `GET /service_types` — Zwraca listę dostępnych typów usług.
- `POST /service_types/search` — Zaawansowane pobieranie dostępnych typów usług.
- `GET /service_types/{_id}` — Zwraca pojedynczy typ usługi.
- `GET /service_types/{_id}/extended_setting_formats` — Zwraca listę specyfikacji formatów zaawansowanych ustawień usług wybranego typu.
- `GET /services/{_id}/type` — Zwraca typ wybranej usługi.

## services

- `GET /accounts/{_id}/services` — Zwraca listę usług należących do konta.
- `GET /groups/{_id}/services` — Zwraca listę usług należących do grupy.
- `GET /scheduled_suspensions/{_id}/service` **[DEPRECATED]** — Zwraca usługę, której dotyczy wybrane zawieszenie.
- `GET /services` — Zwraca listę usług.
- `POST /services` — Tworzy nową usługę.
- `POST /services/bulk_activate` — Włącza monitoring wielu usług jednocześnie.
- `POST /services/bulk_change_group` — Przenosi wiele usług do grupy jednocześnie.
- `POST /services/bulk_delete` — Usuwa wiele usług jednocześnie.
- `POST /services/bulk_force_analysis` — Wymusza wykonanie analizy działania wielu usług jednocześnie poza zaplanowanym harmonogramem.
- `POST /services/bulk_pause` — Wyłącza monitoring wielu usług jednocześnie.
- `POST /services/search` — Zaawansowane pobieranie usług.
- `DELETE /services/{_id}` — Usuwa usługę.
- `GET /services/{_id}` — Zwraca pojedynczą usługę.
- `PUT /services/{_id}` — Modyfikuje ustawienia usługi.
- `POST /services/{_id}/archive` — Archiwizuje usługę.
- `GET /services/{_id}/escalation_statuses` — Zwraca listę aktualnych statusów eskalacji dla usługi.
- `GET /services/{_id}/event_statuses` — Zwraca listę aktualnych statusów zdarzeń dla usługi.
- `POST /services/{_id}/force_analysis` — Wymusza wykonanie analizy działania usługi poza zaplanowanym harmonogramem.
- `POST /services/{_id}/force_custom_analysis` — Wymusza wykonanie analizy działania usługi poza zaplanowanym harmonogramem, z niestandardowymi ustawieniami.
- `GET /services/{_id}/notification_channels` — Zwraca listę kanałów powiadomień włączonych dla usługi.
- `GET /services/{_id}/notification_conditions` — Zwraca listę warunków powiadomień włączonych dla usługi.
- `GET /services/{_id}/notification_mode` — Zwraca tryb powiadomień usługi.
- `GET /services/{_id}/recovery_notification_mode` — Zwraca tryb powiadomień o końcu awarii dla usługi.
- `POST /services/{_id}/restore` — Przywraca zarchiwizowaną usługę.
- `POST /services/{_id}/schedule_custom_analysis` — Wymusza wykonanie analizy działania usługi poza zaplanowanym harmonogramem, z niestandardowymi ustawieniami.
- `POST /services/{_id}/send_test_sms_notification` — Wysyła testowe powiadomienie SMS na numery telefonu przypisane do usługi.
- `GET /suspensions/{_id}/service` — Zwraca usługę, której dotyczy wybrane zawieszenie.
- `GET /weekly_suspensions/{_id}/service` — Zwraca usługę, której dotyczy wybrane tygodniowe zawieszenie.

## sessions

- `GET /sessions` — Zwraca listę sesji API.
- `POST /sessions` — Tworzy nową sesję API.
- `POST /sessions/search` — Zaawansowane pobieranie sesji API.
- `DELETE /sessions/{_id}` — Usuwa sesję API.
- `GET /sessions/{_id}` — Zwraca pojedynczą sesję API.
- `POST /sessions/{_id}/2fa/verify` — Weryfikuje sesję za pomocą jednorazowego kodu 2FA.

## suspensions

- `GET /services/{_id}/suspensions` — Zwraca listę zaplanowanych zawieszeń dla usługi.
- `GET /services/{_id}/weekly_suspensions` — Zwraca listę tygodniowych zawieszeń dla usługi.
- `GET /suspensions` — Zwraca listę zaplanowanych zawieszeń.
- `POST /suspensions` — Tworzy nowe zawieszenie.
- `POST /suspensions/bulk_create` — Dodaje zaplanowane zawieszenie do wielu usług jednocześnie.
- `POST /suspensions/bulk_delete` — Usuwa wiele zawieszeń jednocześnie.
- `POST /suspensions/search` — Zaawansowane pobieranie zaplanowanych zawieszeń.
- `DELETE /suspensions/{_id}` — Usuwa zaplanowane zawieszenie.
- `GET /suspensions/{_id}` — Zwraca pojedyncze zaplanowane zawieszenie.
- `PUT /suspensions/{_id}` — Modyfikuje zaplanowane zawieszenie.
- `GET /weekly_suspensions` — Zwraca listę tygodniowych zawieszeń.
- `POST /weekly_suspensions` — Tworzy nowe tygodniowe zawieszenie.
- `POST /weekly_suspensions/search` — Zaawansowane pobieranie tygodniowych zawieszeń.
- `DELETE /weekly_suspensions/{_id}` — Usuwa tygodniowe zawieszenie.
- `GET /weekly_suspensions/{_id}` — Zwraca pojedyncze tygodniowe zawieszenie.
- `PUT /weekly_suspensions/{_id}` — Modyfikuje tygodniowe zawieszenie.
- `GET /weekly_suspensions/{_id}/permissions` — Zwraca uprawnienia do wybranego zawieszenia.

## system

- `POST /system/dns_diagnostics` — Tworzy nowe zapytanie diagnostyczne DNS.
- `DELETE /system/dns_diagnostics/{_id}` — Usuwa zapytanie diagnostyczne DNS.
- `GET /system/dns_diagnostics/{_id}` — Zwraca pojedyncze zapytanie diagnostyczne DNS.
- `GET /system/my_ip_address` — Zwraca adres IP klienta API.
- `POST /system/reset_password/confirm` — Potwierdza reset hasła do konta.
- `POST /system/reset_password/initialize` — Generuje żądanie resetu hasła do konta.
- `GET /system/stats` — Zwraca aktualne statystyki systemu Monit24.pl.
- `POST /system/traceroute_diagnostics` — Tworzy nowe zapytanie diagnostyczne traceroute.
- `DELETE /system/traceroute_diagnostics/{_id}` — Usuwa zapytanie diagnostyczne traceroute.
- `GET /system/traceroute_diagnostics/{_id}` — Zwraca pojedyncze zapytanie diagnostyczne traceroute.

## templates

- `GET /accounts/{_id}/templates` — Zwraca listę szablonów powiadomień należących do konta.
- `GET /contacts/{_id}/custom_templates` — Zwraca listę własnych szablonów dla kontaktu.
- `DELETE /contacts/{_id}/custom_templates/{_channel_id}` — Usuwa przypisanie własnego szablonu do kontaktu.
- `GET /contacts/{_id}/custom_templates/{_channel_id}` — Zwraca własny szablon dla określonego kontaktu i kanału.
- `PUT /contacts/{_id}/custom_templates/{_channel_id}` — Tworzy lub modyfikuje przypisanie własnego szablonu do kontaktu.
- `GET /contacts/{_id}/custom_templates/{_channel_id}/permissions` — Zwraca uprawnienia do przypisania własnego szablonu do kontaktu.
- `GET /default_templates` — Zwraca listę domyślnych szablonów.
- `POST /default_templates/search` — Zaawansowane pobieranie domyślnych szablonów.
- `GET /default_templates/{_channel_id}/{_language_id}` — Zwraca pojedynczy domyślny szablon.
- `GET /templates` — Zwraca listę szablonów.
- `POST /templates` — Tworzy nowy szablon.
- `POST /templates/search` — Zaawansowane pobieranie szablonów.
- `DELETE /templates/{_id}` — Usuwa szablon.
- `GET /templates/{_id}` — Zwraca pojedynczy szablon.
- `PUT /templates/{_id}` — Modyfikuje szablon.

## user_data

- `GET /user_data/{_id}` — Zwraca dane wybranego użytkownika.
- `PUT /user_data/{_id}` — Uaktualnia dane wybranego użytkownika.
- `DELETE /user_data/{_id}/settings/{_key}` — Usuwa wybrane ustawienie użytkownika.
- `GET /user_data/{_id}/settings/{_key}` — Zwraca wartość wybranego ustawienia użytkownika.
- `PUT /user_data/{_id}/settings/{_key}` — Aktualizuje lub tworzy wartość wybranego ustawienia użytkownika.

