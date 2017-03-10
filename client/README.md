# Client README

## Commands
*Tänk på att man bara kan använda en flagga åt gången.*

### check_ram [-total]
check_ram kollar RAM användningen på klienten.

**-total**

Om -total används så får man klienten's totala RAM, annars får man bara klientens RAM användning.

### check_disc [-total]
check_disc kollar partitions användningen på klienten.

**-total**

Om -total används så får man klienten's totala partition för alla partioner, annars får man bara klientens användning för alla partioner.

### check_cpu [-total] [-percent=[true|false]]
check_cpu kollar CPU användningen på klienten.

**-info**

Om -info används så får man CPU information från klienten, CPU model, family etc.

**-percent=[true|false]**
Om -percent är true så får man tillbaka procent användning per cpu, annars får man tillbaka den totala procent användningen, tänk på att man då kan få tillbaka mer än 100%.

*Default är false*

### update -url=(URL)
update uppdaterar klienten till senaste versionen.

**-url=(URL)**

-url specifierar URLn till version APIn för att man ska uppdatera klienten.