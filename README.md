# App Versions Slack Slash
An Slack Slash Command for [AppVersions](https://appversions.vercel.app/).

## Slack Slash commands
Add the given appIds to the lookup table:
```
/command add [android|ios] appId0 appId1 ...
```

Remove the given appIds from the lookup table
```
/command remove [android|ios] appId0 appId1 ...
```

Get the current (shorted) URL from the saved appIds to [AppVersions](https://appversions.vercel.app/)
```
/command get
```

See the app versions information directly in slack
```
/command lookup [android|ios] appId0 appId1 ...
```

## How to install
### Prerequisite
* [Vercel](https://vercel.com/) account
* Firebase project
* (Obviously) an Slack workspace

### Vercel
Run the vercel cli (`vercel`) to create an project on vercel.
After that we're able to set environment variables for this project.

We have to set the following **environment variables**:

**`FIREBASE_SERVICE_ACCOUNT`**
</br>
The value of the Firebase Service Account for the Firebase Admin SDK.
You can get one under your Firebase project settings -> Sevice Account -> Firebase Admin SDK -> Generate new private key

The content of the `*.json` file has to be placed to this environment variable.

**`FIREBASE_DYNAMIC_LINKS_DOMAIN`**
</br>
As we shorten the [AppVersions](https://appversions.vercel.app/) URL, you have to enable Firebase Dynamic Links and create a new (free) domain for it. You can literally choose any subdomain name of `.page.link`.

Place the URL (e.g. `https://appversions.page.link`) to this environment variable.

**`FIREBASE_DYNAMIC_LINKS_API_KEY`**
</br>
As the Firebase Admin SDK doesn't support dynamic links we have to call their API.
For this we need the "Web API Key" from the Firebase project.
The key can be found under the project settings.

**Note**: In case there is no key visible, then enable Firebase Authentication (but don't enable any Authentication Providers). This will generate the Web API Key for us.

Place the API Key in this environment variable.

**`SLACK_SIGNING_SECRET`**
</br>
Because we are verifying the request, according to the [Slack documentation](https://api.slack.com/authentication/verifying-requests-from-slack), we have to set the Slack Signing Secret as a environment variable.

The Signing Secret can be found under the Slack App "Basic Information" menu under "App Credentials".

Simply store this Key in this environment variable.

**Redeploy the project**
</br>
After we set up all the environment variables we have to redeploy the vercel project.
Otherwise the environment variables aren't visible.

### Firebase

#### Firestore
Beside of the stuff we did already to set up vercel, we also have to enable Firebase Firestore.

Simply click on the menu entry and `Create Database`.
You can use the following rules for the Database.
Those rules make the DB private to everyone, except of the Admin SDK:
```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if false
    }
  }
}
```

#### Dynamic Links
We also have to set up an "Allowlist URL pattern".

The following URL should be added to that list. All others (in case there are one) can be removed:
```
^https://appversions\.vercel\.app\?.*$
```

See also for more information [this guide](https://support.google.com/firebase/answer/9021429?hl=en).


### Slack
We have to create a new Slack App, which contains nothing else than an Slash Command.

You can do this by visiting https://api.slack.com/apps -> Create New App -> Slash Commands -> Create New Command

The following options should be set:

| Option | Value |
|---|---|
| Command | What you want. Recommendet `/appversions` |
| Request Url | https://[SubDomain].vercel.app/api/slashcommand |
| Short description | See and manage the current versions of our mobile applications |
| Usage Hint | [add remove get] [android\|ios] appId |

Click **Save**.

Now go back to the App and click "Install to Workspace" and you're done 🎉.

See also the [official Slack documentation](https://api.slack.com/interactivity/slash-commands#creating_commands) for it.
