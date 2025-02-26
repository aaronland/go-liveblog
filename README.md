# go-guardian-liveblog

Go package for watching a variety of live blog events and reading them aloud using the operating system's text-to-speech APIs.

## Important

This is MacOS specific right now.

## Tools

```
$> make cli
go build -mod vendor -ldflags="-s -w" -o bin/follow cmd/follow/main.go
```

### follow

Parse one or more "live blog" URLs and read them aloud.

```
$> ./bin/follow -h
Parse one or more "live blog" URLs and read them aloud.
Usage:
	 ./bin/follow [options] url(N) url(N)
  -delay int
    	The number of seconds to wait before fetching new updates (default 30)
  -read-all
    	If true read all previous posts (written before following has begun)
  -verbose
    	Enable verbose (debug) logging.
```

The `follow` tool will keep a local cache of posts its already seen (and read) for the duration it is run.

#### Example

##### The Guardian and Le Monde

```
$> ./bin/follow \
	https://www.lemonde.fr/sport/live/2024/08/07/direct-volley-ball-france-italie-suivez-le-match-des-demi-finales-du-tournoi-masculin-aux-jo-2024_6272013_3242.html \
	https://www.theguardian.com/us-news/live/2024/aug/07/kamala-harris-tim-walz-vp-election-campaign-updates \
	https://www.theguardian.com/sport/live/2024/aug/07/paris-2024-olympics-day-12-live-updates-today-schedule-events-athletics-cycling-golf-diving

2024/08/07 11:10:04 INFO           Le speaker du match n'est-il pas le speaker de Roland Garros ?    Rafa            Vous voulez parler de Marc Maury, aussi connu des spectateurs assidus des meetins d’athlétisme ? C’est possible, mais difficile de vous le confirmer devant notre télévision. Une chose est sûre : sa compétence de chauffeur d’arène n’est plus à prouver.                     

2024/08/07 11:10:23 INFO    Aie aie aie, le service sans sel de Barthélémy Chinenyeze…            Flashé à 7 km/h (ok, on exagère un peu) et directement dans le filet. Muscle ton jeu Barthélémy.                     

2024/08/07 11:10:35 INFO Here is another video of the crowd waiting in anticipation of Kamala Harris and Tim Walz in Eau Claire, Wisconsin, where the duo is set to take the stage at around 2.30pm ET:

2024/08/07 11:10:46 INFO Also in the second race is Letsile Tebogo of Botswana. Real talk, this isn’t the strongest field we’ve ever had, but he might be Lyles’ closest challenger; remember, Lyles hasn’t lost in 26 races, – since the final in Tokyo.
```

##### La Presse

```
$> ./bin/follow \
	-read-all \
	https://www.lapresse.ca/sports/hockey/2025-02-25/1re-periode/hurricanes-0-canadien-0.php

2025/02/25 17:06:41 INFO Brent Burns atteint Caufield au visage, il est à son tour chassé. Le Canadien jouera à 5 contre 3 pendant 1 min 17 s
2025/02/25 17:06:49 INFO 
2025/02/25 17:06:49 INFO Le Canadien retourne en avantage numérique. Gostisbehere fait trébucher Dvorak, mais c'est Gallagher qui était à l'origine de la séquence
2025/02/25 17:06:58 INFO Assez brouillon pour le Canadien en 2e jusqu'ici
2025/02/25 17:07:02 INFO Un bel arrêt de Montembeault contre Jarvis dès le départ
2025/02/25 17:07:05 INFO C'est reparti pour la 2e période et de retour à 5 contre 5


2025/02/25 17:07:39 INFO Beaucoup de passes, mais pas de tir. Le triangle défensif ne permet aucune passe dangereuse vers Laine
2025/02/25 17:08:10 INFO Gostisbehere de retour. Le Canadien joue maintenant à 5 contre 4 pendant 35 secondes
2025/02/25 17:09:39 INFO De retour à 5 contre 5. Des décisions douteuses de Laine et Newhook ont fait mal au Canadien. Et les Hurricanes ont montré pourquoi ils sont 1ers dans la LNH en désavantage
```

_Note: User comments (and replies) are not supported yet._