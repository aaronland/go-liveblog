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

#### Example

```
> ./bin/follow \
	https://www.lemonde.fr/sport/live/2024/08/07/direct-volley-ball-france-italie-suivez-le-match-des-demi-finales-du-tournoi-masculin-aux-jo-2024_6272013_3242.html \
	https://www.theguardian.com/us-news/live/2024/aug/07/kamala-harris-tim-walz-vp-election-campaign-updates \
	https://www.theguardian.com/sport/live/2024/aug/07/paris-2024-olympics-day-12-live-updates-today-schedule-events-athletics-cycling-golf-diving

2024/08/07 11:10:04 INFO           Le speaker du match n'est-il pas le speaker de Roland Garros ?    Rafa            Vous voulez parler de Marc Maury, aussi connu des spectateurs assidus des meetins d’athlétisme ? C’est possible, mais difficile de vous le confirmer devant notre télévision. Une chose est sûre : sa compétence de chauffeur d’arène n’est plus à prouver.                     

2024/08/07 11:10:23 INFO    Aie aie aie, le service sans sel de Barthélémy Chinenyeze…            Flashé à 7 km/h (ok, on exagère un peu) et directement dans le filet. Muscle ton jeu Barthélémy.                     

2024/08/07 11:10:35 INFO Here is another video of the crowd waiting in anticipation of Kamala Harris and Tim Walz in Eau Claire, Wisconsin, where the duo is set to take the stage at around 2.30pm ET:

2024/08/07 11:10:46 INFO Also in the second race is Letsile Tebogo of Botswana. Real talk, this isn’t the strongest field we’ve ever had, but he might be Lyles’ closest challenger; remember, Lyles hasn’t lost in 26 races, – since the final in Tokyo.
```