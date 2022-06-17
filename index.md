## The Versātilis Project

Contact [Micah Sherr](https://micahsherr.com) for more information.


### Description

Internet freedom worldwide is becoming increasingly more
restricted[^1] as new censorship practices are developed and deployed
by nation-states seeking to control the flow of information for
political, social, and/or economic gain[^2][^3].  As researchers
continue to better understand the various mechanisms of
censorship[^4], privacy advocates have simultaneously developed
censorship-resistant systems (CRSes) that attempt to evade censorship
systems and allow their users to have unfettered access to the
Internet.

Unfortunately, this has led to an arms-race in which the status quo
significantly favors the censor who has far more human, computational,
and financial resources than those who wish to evade censorship.
Although CRSes have become more advanced, Internet censorship remains
widespread with information access controls in place in at least 60
countries[^5].

This project aims to develop algorithms, communication protocols, and
implementations that achieve censorship resistance as a "side-effect"
of providing other desirable Internet communication features.  This
project takes a radically different approach than prior work towards
tackling censorship resistance: we take the position that rather than
existing as a standalone system, censorship-resistance should be a
characteristic of a widely fielded and general-purpose communication
platform.  We posit that it is more difficult for a censor to block a
ubiquitous and widely-used communication protocol than a niche
application designed solely to circumvent censorship. Our goal is to
avoid the censor vs. anti-censorship arms-race by instrumenting a
reliable and high-performance communication platform that we hope will
be (1) widely deployed, (2) not used primarily as an anti-censorship
apparatus, and (3) inherently difficult to surveil and block as a
natural consequence of its design.

Paradoxically, to be effective as an anti-censorship technology, such
a communication platform should achieve widespread adoption for
purposes unrelated to censorship circumvention.  If the primary
purpose of the platform is censorship circumvention, the cost to the
adversary of barring access to protocols built using the platform is
low.  However, if the platform is also regularly used for business and
commerce, blocking an otherwise useful tool that has widespread
adoption may be too politically and economically costly for a censor.
To this end, the platform must both encourage general purpose usage
and be competitive with existing methods of communication.


### Project Personnel

* [Micah Sherr](https://micahsherr.com), Professor, Department of
Computer Science, Georgetown University.

### Papers

Tan, Henry, and Micah Sherr. "Censorship Resistance as a Side-Effect."
In International Workshop on Security Protocols, 2014.


### Funding

This project is partially funded through the
[McCourt Institute](https://mccourtinstitute.org/).  The opinions and
findings expressed in the products of this project are strictly those
of the author(s) and do not necessarily those of any employer or
funding agency.


---

[^1]: See Freedom House. 2021. Freedom of Expression. https://freedomhouse.org.

[^2]: Arian Akhavan Niaki, Shinyoung Cho, Zachary Weinberg, Nguyen Phong Hoang, Abbas Razaghpanah, Nicolas Christin, and Phillipa Gill. 2020. ICLab: A Global, Longitudinal Internet Censorship Measurement Platform. In IEEE Symposium on Security and Privacy (SP).

[^3]: Ram Sundara Raman, Adrian Stoll, Jakub Dalek, Reethika Ramesh, Will Scott, and Roya Ensafi. 2020. Measuring the Deployment of Network Censorship Filters at Global Scale. In Proceedings of the Network and Distributed System Security Symposium (NDSS).

[^4]: M. C. Tschantz, S. Afroz, Anonymous, and V. Paxson. 2016. SoK: Towards Grounding Censorship Circumvention in Empiricism.  In IEEE Symposium on Security and Privacy (SP).

[^5]: Niaki, *op. cit.*
