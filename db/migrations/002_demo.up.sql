insert into collections (
	id,
	name,
	description
) values (
	'630cfd08-924c-49fa-b5b2-2c81e979829d',
	'Demo collection',
	'Some cards to get a feel of the app'
);

insert into decks (
	id,
	collection_id,
	name,
	description
) values (
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'630cfd08-924c-49fa-b5b2-2c81e979829d',
	'Refresher demo',
	'Just a simple quiz'
);

insert into cards (
	id,
	deck_id,
	content
) values (
	'83f63b6e-7a61-49f6-9fd6-78bbe367088c',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"theme":{"card":{"fill_color":"var(--app-theme-sapphire)","mask_color":"var(--app-theme-snow-white)"},"interactives":{"fill_color":"var(--app-theme-snow-white)","mask_color":"var(--app-theme-midnight)"}},"content":[{"type":"title","content":"Geography 101"},{"type":"textbox","content":[{"type":"text","content":"What''s the current capital of "},{"type":"text","content":"Germany?","theme":{"highlight":{"text_color":"white","fill_color":"var(--app-theme-spooky-orange)"},"bold":true,"decoration":"underline"}}]},{"type":"poll","is_quiz":true,"content":[{"value":"Cologne"},{"value":"Berlin","is_answer":true},{"value":"Bonn"}]}]},"back":{"theme":{"card":{"outline_color":"var(--app-theme-sapphire)"}},"content":[{"type":"textbox","content":[{"type":"text","content":"Berlin has been the capital of Germany again since October 3, 1990"}]}]}}'
), (
	'b2e7070b-1b3e-4aa3-8fab-af7940a47e71',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"theme":{"card":{"fill_color":"var(--app-theme-spooky-orange)","mask_color":"var(--app-theme-snow-white)"},"interactives":{"fill_color":"var(--app-theme-snow-white)","mask_color":"var(--app-theme-spooky-orange)"}},"content":[{"type":"title","content":"Turismus 101"},{"type":"textbox","content":[{"type":"text","content":"In which country it is considered a grave offence to order pizza with pineapple?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Spain"},{"value":"Italy","is_answer":true},{"value":"USA"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"Itallians will kick the shit out of you if you try to deface their national dish by putting pineapple on it"}]}]}}'
), (
	'decf636d-15ca-404b-90a7-d0d6febf948e',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Chemistry 101"},{"type":"textbox","content":[{"type":"text","content":"What is the chemical symbol for Gold?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Ag"},{"value":"Au","is_answer":true},{"value":"Fe"},{"value":"Pb"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"The symbol Au comes from the Latin word for gold, ''aurum''."}]}]}}'
), (
	'b700001d-2186-4611-8985-58b7a1e2e4b0',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Astronomy"},{"type":"textbox","content":[{"type":"text","content":"Which planet is known as the ''Red Planet''?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Venus"},{"value":"Mars","is_answer":true},{"value":"Jupiter"},{"value":"Saturn"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"Mars appears red due to the iron oxide (rust) on its surface."}]}]}}'
), (
	'09eac375-2437-4c68-94be-e9cdc220cab6',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"History"},{"type":"textbox","content":[{"type":"text","content":"In what year did the French Revolution begin?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"1776"},{"value":"1789","is_answer":true},{"value":"1812"},{"value":"1492"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"The revolution began in 1789 with the Storming of the Bastille."}]}]}}'
), (
	'024ebbeb-187d-42d8-9d6b-0ca1313388a0',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Biology"},{"type":"textbox","content":[{"type":"text","content":"Which organelle is known as the ''powerhouse of the cell''?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Nucleus"},{"value":"Ribosome"},{"value":"Mitochondria","is_answer":true},{"value":"Golgi Apparatus"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"Mitochondria generate most of the cell''s supply of adenosine triphosphate (ATP), used as a source of chemical energy."}]}]}}'
), (
	'7f96de51-bb05-42ed-a675-5fc7b844ff95',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Geography"},{"type":"textbox","content":[{"type":"text","content":"Which is the longest river in the world?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Amazon"},{"value":"Nile","is_answer":true},{"value":"Yangtze"},{"value":"Mississippi"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"The Nile is traditionally considered the longest river in the world, stretching about 6,650 kilometers."}]}]}}'
), (
	'53d83e97-2ceb-4bc0-ad89-1e8f8e818a2c',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Literature"},{"type":"textbox","content":[{"type":"text","content":"Who wrote the dystopian novel ''1984''?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Aldous Huxley"},{"value":"George Orwell","is_answer":true},{"value":"Ray Bradbury"},{"value":"Ernest Hemingway"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"George Orwell published ''1984'' in 1949 as a warning against totalitarianism."}]}]}}'
), (
	'59518990-1e17-4d53-9feb-fbce67391a57',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Fine Arts"},{"type":"textbox","content":[{"type":"text","content":"Who painted the ''Mona Lisa''?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"Vincent van Gogh"},{"value":"Leonardo da Vinci","is_answer":true},{"value":"Pablo Picasso"},{"value":"Claude Monet"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"Leonardo da Vinci painted the Mona Lisa in the early 16th century during the Italian Renaissance."}]}]}}'
), (
	'94fe1e9b-e2ab-41a3-95a3-038989e6d49b',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Technology"},{"type":"textbox","content":[{"type":"text","content":"What does ''HTTP'' stand for?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"HyperText Transfer Protocol","is_answer":true},{"value":"High Transfer Text Process"},{"value":"Hyperlink Textual Transfer"},{"value":"Home Tool Transfer Process"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"HTTP is the foundation of data communication for the World Wide Web."}]}]}}'
), (
	'1686ba9d-fff7-4b92-8569-690d2f6cf302',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Music"},{"type":"textbox","content":[{"type":"text","content":"How many symphonies did Ludwig van Beethoven complete?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"5"},{"value":"7"},{"value":"9","is_answer":true},{"value":"12"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"Beethoven''s 9th Symphony, which includes the ''Ode to Joy'', was his final completed symphony."}]}]}}'
), (
	'b5ec10eb-edaf-4191-be6c-ef2473ad6845',
	'33511bb5-f9d7-4795-9b3a-c1479378c27b',
	'{"front":{"content":[{"type":"title","content":"Mathematics"},{"type":"textbox","content":[{"type":"text","content":"What is the value of Pi ($pi$) to two decimal places?"}]},{"type":"poll","is_quiz":true,"content":[{"value":"3.12"},{"value":"3.14","is_answer":true},{"value":"3.16"},{"value":"3.18"}]}]},"back":{"content":[{"type":"textbox","content":[{"type":"text","content":"$pi$ is the ratio of a circle''s circumference to its diameter, approximately equal to 3.14159."}]}]}}'
);
