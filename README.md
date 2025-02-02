# tp_gobuster
# Outil Gobuster-like en Go

## Description
Ce projet implémente un outil en Go similaire à **Gobuster**, utilisé pour identifier les fichiers et répertoires cachés sur un serveur web en effectuant des requêtes HTTP avec une liste de mots (dictionnaire). Ce type d'outil est couramment utilisé lors de la phase de reconnaissance dans des tests d'intrusion.

## Fonctionnalités
L'outil fonctionne de manière simple :
- Utilisation d'un dictionnaire de mots pour tester différentes URL sur un serveur cible.
- Affichage des réponses HTTP pour chaque URL testée.
- Possibilité d'exécuter plusieurs requêtes en parallèle pour améliorer les performances.

## Usage
Le programme peut être exécuté avec les différentes options suivantes :

### Flags
- `-d` : **Chemin vers le fichier dictionnaire** (fichier texte contenant une liste de mots à tester).
- `-q` : **Mode silencieux**. Affiche uniquement les résultats HTTP 200 (c’est-à-dire les URL valides).
- `-t` : **Cible à scanner** (l'URL du serveur à tester, peut inclure un port).
- `-w` : **Nombre de workers** (nombre de threads à exécuter en parallèle, la valeur par défaut est 1).

### Exemple d'utilisation
#### Exécution simple avec un fichier dictionnaire :
```bash
go run main.go -w 10 -d /usr/share/dict/words -t http://localhost:8080

Contribution
Vous pouvez cloner ce dépôt et y apporter des améliorations. Merci de respecter les bonnes pratiques de développement et de créer des pull requests pour toute contribution.

Avertissements
Assurez-vous que le serveur cible est opérationnel avant de lancer un scan.
L'outil est destiné à être utilisé dans des environnements légaux, notamment pour des tests de pénétration autorisés.

