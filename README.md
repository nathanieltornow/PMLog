# PMLog

**Requirements**
- [Go](https://golang.org/) (version >= 1.15). If you're using Linux, you can use `./scripts/install_go.sh && source /etc/profile`.
- [Make](https://www.gnu.org/software/make/)


### Start a sequencer

1. Compile
   ```shell
   make seq
   ```

2. Start a sequencer
   ```shell
   ./seq -IP <ip-address> -color <colorID> [-parIP <parent-IP-address>] [-root] [-leader] 
   ```
   ```text
   Usage of ./seq:
   -IP string
   The IP on which the sequencer listens
   -color int
   The color which the sequencer represents
   -leader
   If the sequencer is a leader and can therefore reply to requests
   -parIP string
   The IP of the parent-sequencer
   -root
   If the sequencer is the root-sequencer
   ```
