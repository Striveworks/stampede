[![Striveworks][striveworks-shield]][license-url]
[![MIT License][license-shield]][license-url]



<!-- PROJECT LOGO -->
<br />
<p align="center">
  <a href="https://github.com/othneildrew/Best-README-Template">
    <img src="stampede.png" alt="Logo" width="600" height="300">
  </a>

  <h3 align="center">Stampede</h3>

  <p align="center">
    Bootstrap Microk8s clusters
    <br />
    <a href=""><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="">View Demo</a>
    ·
    <a href="">Report Bug</a>
    ·
    <a href="">Request Feature</a>
  </p>
</p>



<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project
Stampede is meant to make bootstrapping Microk8s clusters seamless. It uses a simple election protocol to elect a leader and then followers will follow. The leader bootstraps the cluster and deals out join tokens to any followers. All communication is handled via a specified multicast channel.

This project is inspired by Chic-fil-A's [Highlander](https://medium.com/@cfatechblog/bare-metal-k8s-clustering-at-chick-fil-a-scale-7b0607bd3541)

### Built With


* [Go](https://golang.org/)
* [Microk8s](https://microk8s.io/)


<!-- GETTING STARTED -->
## Getting Started

This can be run on any Ubuntu distribution. The install script will create systemd service and run it. Optionally, there is a Vagrant setup that can be used to spin up 3 VMs and bootstrap them into a microk8s cluster.

### Prerequisites

Currently, this project only offers support for Ubuntu distributions

### Installation

1. `make install`


### Testing

1. `make test`


### Generating documentation

1. `make docs`

<!-- ROADMAP -->
## Roadmap

See the [open issues]() for a list of proposed features (and known issues).



<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE` for more information.



<!-- CONTACT -->
## Contact

[Striveworks](striveworks.us)


[striveworks-shield]: https://img.shields.io/badge/BUILT%20BY-STRIVEWORKS-orange?style=for-the-badge
[license-shield]: https://img.shields.io/github/license/othneildrew/Best-README-Template.svg?style=for-the-badge
[license-url]: https://github.com/othneildrew/Best-README-Template/blob/master/LICENSE.txt
