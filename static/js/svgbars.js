
var width = 960;
var height = 500;

var y = d3.scale.linear().range([height, 0]);

var barHeight = 20;

var chart = d3.select(".chart")
    .attr("width", width)
    .attr("height", height);

d3.tsv("static/data/charts.tsv", type, function(error, data) {
    y.domain([0, d3.max(data, function(d){return d.value;})]);

    var barWidth = width / data.length

    var bar = chart.selectAll("g")
        .data(data)
        .enter().append("g")
        .attr("transform", function(d, i){ return "translate(" + i * barWidth + ",0)"; });

    bar.append("rect")
        .attr("y", function(d) { return y(d.value); })
        .attr("height", function(d){ return height - y(d.value); })
        .attr("width", barWidth -1);

    bar.append("text")
        .attr("y", function(d) { return x(d.value) + 3; })
        .attr("x", barWidth/2)
        .attr("dy", ".75em")
        .text(function(d) { return d.value; });
});

function type(d) {
    d.value = +d.value; //coerce value to integer
    return d;
}
