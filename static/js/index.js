// impression
function impCount(diff, divide) {
    return function(index, ele) {
        d = new Date();
        unix = d / 1000;
        ele.innerText = parseInt((unix - diff) * (ele.dataset.seed/divide));
    }
}

$(".impression").each(impCount(1508000000, 1));

$(".down-vote-count").each(impCount(1508175000, 85));

$(".up-vote-count").each(impCount(1508150000, 45));
