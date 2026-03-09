package nlp

import "strings"

// stopwords is built at init from the ISO stopword lists plus custom additions.
var stopwords map[string]struct{}

func init() {
	stopwords = make(map[string]struct{}, 1200)

	// ISO 639 English stopwords (stopwords-iso/stopwords-en)
	addWords(isoEnglish)
	// ISO 639 French stopwords (stopwords-iso/stopwords-fr), accent-stripped
	addWords(isoFrench)
	// Custom job-posting boilerplate and city names
	addWords(customStopwords)

	// Remove technical terms that should NOT be stopwords.
	// These are real skills/tools that matter for ATS matching.
	for _, w := range strings.Fields(protectedTerms) {
		delete(stopwords, w)
	}
}

func addWords(block string) {
	for _, w := range strings.Fields(block) {
		if len(w) >= 2 {
			stopwords[w] = struct{}{}
		}
	}
}

func isStopword(token string) bool {
	_, ok := stopwords[token]
	return ok
}

// protectedTerms are technical keywords that must NOT be treated as stopwords
// even if they appear in the ISO lists. Space-separated.
const protectedTerms = `api apis sql nosql css html http https tcp udp grpc json xml yaml csv pdf
rest restful graphql oauth saml sso jwt
react angular svelte nextjs nodejs express jquery redux
docker kubernetes helm istio terraform ansible
aws gcp azure lambda ec2 s3 vpc iam
postgresql postgres mysql mongodb redis elasticsearch cassandra dynamodb firebase sqlite oracle
python java javascript typescript golang rust ruby swift kotlin scala php perl
csharp cpp dotnet aspnet fsharp
git github gitlab bitbucket jenkins circleci
kafka spark hadoop airflow etl
tensorflow pytorch scikit pandas numpy
agile scrum kanban jira
cicd devops sre
linux unix windows macos
nginx apache
blockchain crypto web3 iot
figma sketch adobe
ui ux wcag aria
llm gpt bert rag
springboot spring django flask fastapi express
android ios flutter
regex async websocket mqtt protobuf
microservice microservices serverless
automated automation testing test tests unit integration
ci cd
go react node
solidworks catia autocad ansys fea cfd gdt tolerance machining cnc fmea sixsigma plm pneumatics hydraulics
pcb altium kicad microcontroller plc scada inverter simulink oscilloscope
lean kaizen oee supplychain logistics inventory warehouse procurement forecasting throughput erp sap wms
aerodynamics propulsion avionics aircraft spacecraft satellite orbital arinc do178 do254 verification validation safety`

// isoEnglish contains the stopwords-iso English list (space-separated, lowercase).
// Source: https://github.com/stopwords-iso/stopwords-en
const isoEnglish = `a able about above abroad according accordingly across actually
added adj af affected affecting affects after afterwards again against ago ahead
all allow allows almost alone along alongside already also although always am amid
amidst among amongst amount an and announce another any anybody anyhow anymore anyone
anything anyway anyways anywhere apart apparently appear appreciate appropriate
approximately are area areas aren arent arise around as aside ask asked asking asks
associated at auth available aw away awfully
back backed backing backs backward backwards based basis be became because become becomes
becoming been before beforehand began begin beginning beginnings begins behind being
beings believe below beside besides best better between beyond big bill billion both
bottom brief briefly bring brought bs bt build building business but buy by
call came can cannot cant caption career case cases cause causes certain certainly change
changes challenge challenges challenging clear clearly come comes complete completing
completion computer con concerning consequently consider considering contain containing
contains corresponding could couldn couldnt course create creating cry currently
date day days dear definitely deliver delivering describe described description despite
detail did didn didnt differ different differently difficult directly diverse diversity do
does doesn doesnt doing don done dont doubtful down downed downing downs downwards due during
each early ed effect effort efforts eg eight eighty either eleven else elsewhere empty
enable enabling end ended ending ends enough ensure ensuring entirely especially et even
evenly ever every everybody everyone everything everywhere ex exactly example except expect expected
face faces fact facts fairly far farther felt few fewer fifth fifty fill find finds
fire first five fix focus focused follow followed following follows for forever former
formerly forth forty forward found four free from front full fully further furthered
furthering furthermore furthers
gave general generally get gets getting give given gives giving go goes going gone good
goods got gotten great greater greatest group grouped grouping groups grow growing
had hadn hadnt half happens hardly has hasn hasnt have haven havent having he hed hell
hello help hence her here hereafter hereby herein hers herself hes hi high higher highest
highly him himself his hither home homepage hopefully how however hundred
id ie if ignored ill im immediate immediately importance important improve in inasmuch inc
indeed index indicate indicated indicates information inner inside insofar instead interest
interested interesting interests into invention inward is isn isnt it itd itll its
itself ive
join just
keep keeps kept key kind knew know known knows
large largely last lately later latest latter latterly least length less lest let lets level
like liked likely likewise line little long longer longest look looking looks low lower
made mainly make makes making man manage managing many may maybe me mean means meantime
meanwhile member members men merely might million mine minus miss more moreover most
mostly move mr mrs ms much mug must my myself
name namely near nearly necessarily necessary need needed needing needs neither never
nevertheless new newer newest next nine ninety nobody non none nonetheless noone nor normally
not noted nothing notwithstanding novel now nowhere number numbers
obtain obtained obviously of off offer offering offers often oh ok okay old older oldest
on once one ones only onto open opened opening opens or order ordered ordering orders
org other others otherwise ought our ours ourselves out outside over overall owing own
page part parted particular particularly parting parts past people per perhaps person
place placed places plan planning play please plus point pointed pointing points poorly
positive possible possibly potentially present presented presenting presents presumably
previously primarily probably problem problems process project projects promptly proud
provide provided provides providing put puts
quite
ran range rather rd re readily really reasonably recent recently ref refs regarding regardless
regards related relative relatively relevant report reporting research reserved respectively
respect responsible responsibility result resulted resulting results right room rooms round run running
said same saw say saying says second secondly seconds section see seeing seem seemed seeming
seems seen sees self selves sensible sent serious seriously serve serving set seven seventy several
shall share sharing she shed shell shes short show showed showing shown shows side sides
significant significantly similar similarly since sincere site six sixty skill skills
slight slightly small smaller smallest so some somebody someday somehow someone somethan
something sometime sometimes somewhat somewhere soon sorry specific specifically specified specify
specifying start starting state states still stop strongly sub substantially success
successful successfully such sufficiently suggest sup support supported supporting sure
system
take taken taking tell ten tend tends th than thank thanks that thats the their theirs them
themselves then thence there thereafter thereby therefore therein thereof thereupon these
they theyd theyll theyre theyve thing things think thinks third thirty this thorough
thoroughly those thou though thought thoughts thousand three through throughout thru
thus til till time tip to today together too took top toward towards tried tries truly try
trying turn turned turning turns twelve twenty twice two
under underneath unless unlike unlikely until up upon ups upwards us use used useful
usefully usefulness uses using usually utilize
value values various very via
want wanted wanting wants was wasn wasnt way ways we wed well went were weren werent weve
what whatever when whence whenever where whereafter whereas whereby wherein whereupon
wherever whether which whichever while whilst whither who whoever whole whom whomever
whose why widely will willing wish with within without won wont wonder word words work
worked working works world would wouldn wouldnt
year years yes yet you youd youll young younger youngest your youre yours yourself yourselves`

// isoFrench contains the stopwords-iso French list (accent-stripped, space-separated).
// Source: https://github.com/stopwords-iso/stopwords-fr
// Accents stripped: e for e-accent, a for a-accent, etc.
const isoFrench = `abord absolument afin ah ai aie aient aies ailleurs ainsi ait
allaient allo allons alors anterieur anterieure anterieures apres as assez attendu au
aucun aucune aucuns aujourdhui aupres auquel aura aurai auraient aurais aurait auras
aurez auriez aurions aurons auront aussi autant autre autrefois autrement autres autrui
aux auxquelles auxquels avaient avais avait avant avec avez aviez avions avoir avons ayant
ayez ayons
bah bas basee bat beau beaucoup bien bigre bon boum bravo brrr
car ce ceci cela celle celles celui cependant certain certaine certaines certains certes
ces cet cette ceux chacun chacune chaque cher chers chez chiche chut chere cheres ci cinq
cinquantaine cinquante cinquantieme cinquieme combien comme comment comparable comparables
compris concernant contre couic crac
da dans de debout dedans dehors deja dela depuis dernier derniere derriere des desormais
desquelles desquels dessous dessus deux deuxieme deuxiemement devant devers devra devrait
different differentes differents dire directe directement dit dite dits divers diverse
diverses dix dixhuit dixneuf dixsept dixieme doit doivent donc dont dos douze douzieme
dring droite du duquel durant des debut desormais
effet egale egalement egales eh elle elles en encore enfin entre envers environ es essai
est et etant etc etre eu eue eues euh eurent eus eusse eussent eusses eussiez eussions
eut eux exactement excepte extenso exterieur eumes eut eutes
fais faisaient faisant fait faites facon feront fi flac floc fois font force furent fus
fusse fussent fusses fussiez fussions fut fumes fut futes
gens
ha haut hein hem hep hi ho hola hop hormis hors hou houp hue hui huit huitieme hum hurrah
he helas
ici il ils importe
je jusqu jusque juste
la laisser laquelle las le lequel les lesquelles lesquels leur leurs longtemps lors
lorsque lui la les
ma maint maintenant mais malgre maximale me meme memes merci mes mien mienne miennes
miens mille mince mine minimale moi moindres moins mon mot moyennant multiple multiples
meme memes
na naturel naturelle naturelles ne neanmoins necessaire necessairement neuf neuvieme ni
nombreuses nombreux nommes non nos notamment notre nous nouveau nouveaux nul neanmoins
notre notres
oh ohe olle ole on ont onze onzieme ore ou ouf ouias oust ouste outre ouvert ouverte
ouverts ou
paf pan par parce parfois parle parlent parler parmi parole parseme partant particulier
particuliere particulierement pas passe pendant pense permet personne personnes peu peut
peuvent peux pff pfft pfut pif pire piece plein plouf plupart plus plusieurs plutot
possessif possessifs possible possibles pouah pour pourquoi pourrais pourrait pouvait
prealable precisement premier premiere premierement pres probable probante procedant
proche psitt pu puis puisque pur pure
quand quant quanta quarante quatorze quatre quatrevingt quatrieme quatriemement que quel
quelconque quelle quelles quelquun quelque quelques quels qui quiconque quinze quoi quoique
rare rarement rares relative relativement remarquable rend rendre restant reste restent
restrictif retour revoici revoila rien
sa sacrebleu sait sans sapristi sauf se sein seize selon semblable semblaient semble
semblent sent sept septieme sera serai seraient serais serait seras serez seriez serions
serons seront ses seul seule seulement si sien sienne siennes siens sinon six sixieme soi
soient sois soit soixante sommes son sont sous souvent soyez soyons specifique specifiques
speculatif stop strictement subtiles suffisant suffisante suffit suis suit suivant
suivante suivantes suivants suivre sujet superpose sur surtout
ta tac tandis tant tardive te tel telle tellement telles tels tenant tend tenir tente tes
tic tien tienne tiennes tiens toc toi ton touchant toujours tous tout toute toutefois
toutes treize trente tres trois troisieme troisiemement trop tu
un une unes uniformement unique uniques uns
va vais valeur vas vers vif vifs vingt vivat vive vives vlan voici voie voient voila
voire vont vos votre vous vu ve votre votres
zut ca etaient etais etait etant etat etiez etions ete etee etees etes etes etre`

// customStopwords are job-posting boilerplate, city names, and domain-specific
// noise words that aren't in the ISO lists. Space-separated.
const customStopwords = `role position company team looking seeking responsibilities
requirements qualifications preferred required ability including within across strong
excellent proven experience work working well environment opportunity ideal candidate
applicant apply equal employer benefits salary competitive
knowledge understanding familiar familiarity asset assets availability
career challenge challenging deliver delivering description different
ensure ensuring everyone effort efforts expect expected
information involve involved responsible responsibility
thrive inspire inspiring confirm
act acting enable enabling various continue continued continuing
maximum application applications

poste entreprise equipe recherche recherchons responsabilites competences requises
souhaitees profil candidat postuler salaire avantages environnement annees niveau
recrutement recruter agence talents talent collegues collegue diversite inclusion
inclusif inclusive ensemble participer participation anglais francais bilingue
offre offrons proposons rejoindre rejoignez passionnee passionne dynamique motivee
motive carriere emploi stage stagiaire connaissances connaissance apprentissage
apprendre contenu contenus personnalises personnalise favoriser enrichir disponibles
disponible possedant differentes expertises expertise experiences profils diversifies
points vue titre positif organisation grace permettent permet maitriser metier mode
qualite qualites
bases atout comprehension prerequis prealables secteur activite curiosite fort esprit
rigueur travail completement connexe etudes relever defis supporter croissance
confirmer livrables integrer inspirante respecte meilleures pratiques innovantes
innovante valeur divers partenaires affaires

acces acceder agir action actions assurant assurer assure creer creation cree cours
banque bancaire nationale national back end front developper developpe developpee
deployer deploye fonctions fonction solutions solution technologiques technologique
programmes programme basees basee impact impactant baccalaureat diplome diplomes
type types developpeurs developpeuse assurance resultats resultat processus
procedure procedures gerer gerant tant vient venir supportent supporte concernant
concerne permettre permettant doit doivent presente presenter souhaite souhaiter
capable capacite necessaire necessaires importante essentiels essentiel essentielle
specifique specifiques pertinent pertinente pertinents contribuer repondre ameliorer

montreal toronto vancouver ottawa quebec paris lyon marseille toulouse bordeaux
york san francisco london berlin remote`
