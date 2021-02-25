Chick-fil-a has a great system for managing small k8s clusters at the edge. This project is inspired by their [Highlander](https://medium.com/@cfatechblog/bare-metal-k8s-clustering-at-chick-fil-a-scale-7b0607bd3541) project. Stampede aims to minimize human interaction in cluster setup. It will
magically create a functional kubernetes cluster on startup.

At Striveworks, we use this to quickly bootstrap our on-prem kubernetes clusters.
