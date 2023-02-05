# golibretranslate


<h1>free api translation for golang</h1>

<h2 dir="auto"><a id="user-content-installation" class="anchor" aria-hidden="true" href="#installation"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>Installation</h2>
<div class="highlight highlight-source-shell notranslate position-relative overflow-auto" dir="auto" data-snippet-clipboard-copy-content="go get github.com/kamva/mgm/v3"><pre>go get github.com/antoniomralmeida/golibretranslate</pre></div>




<h2 dir="auto"><a id="user-content-usage" class="anchor" aria-hidden="true" href="#usage"><svg class="octicon octicon-link" viewBox="0 0 16 16" version="1.1" width="16" height="16" aria-hidden="true"><path fill-rule="evenodd" d="M7.775 3.275a.75.75 0 001.06 1.06l1.25-1.25a2 2 0 112.83 2.83l-2.5 2.5a2 2 0 01-2.83 0 .75.75 0 00-1.06 1.06 3.5 3.5 0 004.95 0l2.5-2.5a3.5 3.5 0 00-4.95-4.95l-1.25 1.25zm-4.69 9.64a2 2 0 010-2.83l2.5-2.5a2 2 0 012.83 0 .75.75 0 001.06-1.06 3.5 3.5 0 00-4.95 0l-2.5 2.5a3.5 3.5 0 004.95 4.95l1.25-1.25a.75.75 0 00-1.06-1.06l-1.25 1.25a2 2 0 01-2.83 0z"></path></svg></a>Usage</h2>

<div class="highlight highlight-source-go notranslate position-relative overflow-auto" dir="auto" data-snippet-clipboard-copy-content="import (
   &quot;github.com/kamva/mgm/v3&quot;
   &quot;go.mongodb.org/mongo-driver/mongo/options&quot;
)"><pre><span class="pl-k">import</span> (
   <span class="pl-s">"fmt"</span>
   <span class="pl-s">"github.com/antoniomralmeida/golibretranslate"</span>
)

<span class="pl-k">func</span> <span class="pl-en">main</span>() {
   <span class="pl-c">&emsp;txt, err := golibretranslate.Translate("Bom Dia!", "auto", "en")</span>
   <span class="pl-c">&emsp;if err == nil {</span>
   <span class="pl-c">&emsp;&emsp;fmt.Println(txt)</span>
   <span class="pl-c">&emsp;}</span>
}   
</pre></div>

