window.addEventListener("load", function load(event){

    const sources_el = document.querySelector("#sources");    
    const voices_el = document.querySelector("#voices");
    const messages_el = document.querySelector("#messages");

    var working;
    var voices;
    
    const list_voices = function(){

	voices = window.speechSynthesis.getVoices();
	const count_voices = voices.length;

	if (count_voices == 0){
	    setTimeout(list_voices, 200);
	    return;
	}
	    
	for (var i = 0; i < count_voices; i++) {
	    
	    const opt = document.createElement("option");
	    opt.setAttribute("class", "voice");
	    opt.setAttribute("value", i);
	    
	    if (voices[i].default){
		opt.setAttribute("selected", "selected");
	    }
	    
	    const label = voices[i].name + " (" + voices[i].lang + ")";
	    opt.appendChild(document.createTextNode(label));
	    
	    voices_el.appendChild(opt);
	}

	voices_el.style.display = "block";
    };
    
    const scrub = function(raw) {
	const doc = new DOMParser().parseFromString(raw, 'text/html');
	return doc.body.textContent || "";
    };

    const speak = function(text) {
	
	const utterance = new SpeechSynthesisUtterance(text);
	
	utterance.pitch = 1; 
	utterance.rate = 1;

	const v = parseInt(voices_el.value);
	utterance.voice = voices[v]; 

	window.speechSynthesis.speak(utterance);
    };

    const display = function(text){
	console.log("Display", text);

	const dt = new Date();
	const ts = dt.toLocaleTimeString();

	const time_el = document.createElement("span");
	time_el.setAttribute("class", "timestamp");
	time_el.appendChild(document.createTextNode(ts));

	const msg_el = document.createElement("span");
	msg_el.setAttribute("class", "message");
	msg_el.appendChild(document.createTextNode(text));

	const li_el = document.createElement("li");
	li_el.appendChild(time_el);
	li_el.appendChild(msg_el);

	messages_el.prepend(li_el);	
    };

    const process = function(rsp){
	
	for (const raw of rsp.posts){
	    const p = scrub(raw);
	    speak(p);
	    display(p);
	}
    };

    const get_posts_from_text = function(uri, text){

	working = true;
	
	get_posts(uri, text).then((rsp) => {
	    
	    try {
		process(JSON.parse(rsp));
	    } catch(err){
		console.error("Failed to handle WASM response", err);
	    }

	    working = false;
	    
	}).catch((err) => {
	    console.error("Failed to get posts from WASM", err);
	    working = false;	    
	});
	
    };
    
    const fetch_sources = function(){

	if (working){
	    return;
	}

	working = true;

	if (sources.value == ""){
	    working = false;
	    return;
	}
	
	const uris = sources.value.split("\n")
	
	if (uris.length == 0){
	    working = false;
	    return;
	}

	for (const uri of uris) {

	    if (uri == "http://random.localhost"){
		get_posts_from_text(uri, "");
	    } else {

		console.log("fetch", uri);
		
		fetch(uri).then((rsp) => {
		    get_posts_from_text(uri, rsp);
		}).catch((err) => {
		    display(err);
		    console.error("Failed to fetch", uri, err)
		});
	    }
	}

	working = false;
    };
    
    list_voices();
    
    sfomuseum.golang.wasm.fetch("wasm/get-posts.wasm").then((rsp) => {

	sources.removeAttribute("disabled");
	setInterval(fetch_sources, 1000 * 10);
	
    }).catch((err) => {
	console.error("Failed to load update get sources binary", err);
        return;
    });

});
