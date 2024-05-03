
const addMessage = (msg, owner) => {
  const author = owner === "mine" ? "You" : "ChatGPT";
  const div = document.createElement('div');
  div.className = "message " + owner;
  div.innerHTML = '<div class="author">'+author+':</div><div class="content">'+msg+'</div>';
  document.querySelector('.results').appendChild(div);
};
const disableSubmit = () => { document.querySelector('[type="submit"]').classList.add('disabled'); document.querySelector('.loading').classList.remove('disabled'); };
const enableSubmit = () => {
  document.querySelector('[type="submit"]').classList.remove('disabled');
  document.querySelector('.loading').classList.add('disabled');
  document.querySelector('textarea').classList.remove('disabled');
};
const enableStep2 = () => { document.querySelector('form').classList.remove('step-1'); document.querySelector('form').classList.add('step-2'); };
const addPlayer = videoId => {
  const div = document.createElement('div');
  div.innerHTML = '<iframe src="https://www.youtube.com/embed/'+videoId+'" frameborder="0" allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share" referrerpolicy="strict-origin-when-cross-origin" allowfullscreen></iframe>';
  const iframe = div.querySelector('iframe');
  document.querySelector('form').prepend(iframe);
};
const getVideoId = () => {
  const url = document.querySelector('[name="videoId"]').value.trim();
  if(url.match(/^[^#\&\?]+$/)) {
    return url;
  }
  const regExp = /^.*(youtu\.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/;
  const match = url.match(regExp);
  if (match && match[2].length == 11) {
    return match[2];
  }
  return "";
}

const getQuestion = () => document.querySelector('[name="question"]').value;
const clearQuestion = () => document.querySelector('[name="question"]').value = "";
const sendQuestion = (videoId, question) => {
  if(!videoId) {
    alert("Not a recognizable YouTube video");
    enableSubmit();
    return;
  }
  const isStart = document.querySelectorAll('.step-1').length > 0;
  const url = isStart ? '/start' : '/send';
  //const msg = isStart ? 'Loading Video Summary' : question;
  if(isStart) {
    enableStep2();
    addPlayer(videoId);
  } else {
    addMessage(question, 'mine');
  }
  clearQuestion();
  fetch(url,{
    method: "POST",
    body: JSON.stringify({videoId,question}),
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    }
  })
  .then(r=>{
    if(r.ok) {
      return r.json()
    }
    throw Error(r)
  })
  .then(r=>{
    enableSubmit();
    addMessage(r.choices[0].message.content, '');
  })
  .catch(e=>{
    console.log(e);
    enableSubmit();
  });
}
document.addEventListener('DOMContentLoaded', () => {
  document.querySelector('form').addEventListener('submit', evt => {
    evt.preventDefault();
    disableSubmit();
    sendQuestion(getVideoId(), getQuestion())
    return false;
  });
});
