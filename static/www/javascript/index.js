window.addEventListener("load", function load(event){

    console.log("INIT");
    
    const voices_el = document.querySelector("#voices");
    const messages_el = document.querySelector("#messages");

    var voices;
    
    const default_voice = 22;
    
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
	    
	    if (i == default_voice){
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

	if (v == NaN){
	    v = default_voice;
	}
	
	utterance.voice = voices[v]; 

	window.speechSynthesis.speak(utterance);
    };

    list_voices();
    
    const eventSource = new EventSource('/sse');
    
    eventSource.onmessage = function(event) {
	console.log("New message received:", event.data);

	const data = scrub(event.data);
	speak(data);

	const dt = new Date();
	const ts = dt.toLocaleTimeString();

	const time_el = document.createElement("span");
	time_el.setAttribute("class", "timestamp");
	time_el.appendChild(document.createTextNode(ts));

	const msg_el = document.createElement("span");
	msg_el.setAttribute("class", "message");
	msg_el.appendChild(document.createTextNode(data));

	const li_el = document.createElement("li");
	li_el.appendChild(time_el);
	li_el.appendChild(msg_el);

	messages_el.prepend(li_el);
    };

    eventSource.addEventListener('update', function(event) {
	console.log("Update received:", event.data);
    });
    
    eventSource.onopen = () => {
	console.log("Connection established!");
    };
    
    eventSource.onerror = (error) => {
	console.error("SSE Connection failed or closed.", error);
    };

});
