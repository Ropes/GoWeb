var margin = {top: 20, right:30, left: 30, bottom: 40},
    width = 960 - margin.left - margin.right,
    height = 500 - margin.top - margin.bottom;

var y = d3.scale.linear().range([height, 0]);

var barHeight = 20;

var chart = d3.select(".chart")
    .attr("width", width + margin.left + margin.right)
    .attr("height", height + margin.top + margin.bottom)
    .append("g")
    .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

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
        .attr("y", function(d) { return y(d.value) + 3; })
        .attr("x", barWidth/2)
        .attr("dy", ".75em")
        .text(function(d) { return d.value; });
});

function type(d) {
    d.value = +d.value; //coerce value to integer
    return d;
}
