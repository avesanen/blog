(function(){

	var dropzone = document.getElementById('dropzone');
	var editor = ace.edit("editor");
	var submit = document.getElementById('submit');

	editor.setTheme("ace/theme/chaos");
	editor.getSession().setMode("ace/mode/markdown");
	editor.$blockScrolling = Infinity;
	editor.setOption("showPrintMargin", false);
	editor.setOption("showGutter", false);
	editor.resize();
	editor.focus();

	var upload = function(files) {
		var formData = new FormData();
		for (var i = 0; i < files.length; i++) {
		  formData.append('files', files[i]);
		}

		// now post a new XHR request
		var xhr = new XMLHttpRequest();
		xhr.open('POST', '/upload/');
		xhr.onload = function () {
		  console.log(xhr.responseText);
		  if (xhr.status === 200) {
		    response = JSON.parse(xhr.responseText);
			if (response.status === "ok") {
				console.log(response);
				for (var i = 0; i < response.urls.length; i++) {
					editor.insert("!["+response.urls[i]+"]("+response.urls[i]+")\n")
				}
			}
		  } 
		};
		xhr.send(formData);
	}

	dropzone.ondrop = function(e) {
		e.preventDefault();
		upload(e.dataTransfer.files);
		return false;
	}

	dropzone.ondragover = function() {
		return false;
	}

	dropzone.ondragleave = function() {
		return false;
	}

	submit.onclick = function () {
		console.log("Submit");
		submitEditor(editor.getSession().getValue());
		return false;
	}

	var submitEditor = function(content) {
		var xhr = new XMLHttpRequest();   // new HttpRequest instance 
		xhr.open("POST", "./edit");
		xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
		xhr.onload = function () {
		  console.log(xhr.responseText);
		  if (xhr.status === 200) {
		    response = JSON.parse(xhr.responseText);
			if (response.status === "ok") {
			  window.location.href = response.url;
			}
		  } 
		};
		xhr.send(JSON.stringify({content:content}));
/*
		var formData = new FormData();
		formData.append('content', content);
		var xhr = new XMLHttpRequest();
		xhr.open('POST', './edit');
		xhr.onload = function () {
		  console.log(xhr.responseText);
		  if (xhr.status === 200) {
		    response = JSON.parse(xhr.responseText);
			if (response.status === "ok") {
				window.location.href = response.url;
			}
		  } 
		};
		xhr.send(formData);
		*/
	}

}());