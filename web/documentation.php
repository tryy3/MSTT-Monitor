<div class="col-md-2"></div>

<div class="col-md-8">
    <h1>Commands</h1>
    <p><i>
        () = något värde (true|false|sträng...)<br>
        [] = friviligt
    </i></p>
    <h3>check_ram [-total] [-swap]</h3>
    <p>check_ram kollar RAM användningen på klienten.<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-total </b>Om -total används så får man klienten's totala RAM, annars får man bara klientens RAM användning.<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-swap </b>Om -swap används så får man klienten's swap istället för RAM, man kan använda -total flaggan också.
    </p>

    <h3>check_disc [-total]</h3>
    <p>check_disc kollar partitions användningen på klienten.<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-total </b>Om -total används så får man klienten's totala partition för alla partioner, annars får man bara klientens användning för alla partioner.
    </p>

    <h3>check_cpu [-info]</h3>
    <p>check_cpu kollar CPU användningen på klienten.<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-info </b>Om -info används så får man CPU information från klienten, CPU model, family etc.<br>
    </p>

    <h3>update -url=(URL)</h3>
    <p>update uppdaterar klienten till senaste versionen.<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-url=(URL) </b>-url specifierar URLn till version APIn för att man ska uppdatera klienten.
    </p>

    <h3>uptime -boot</h3>
    <p>uptime skickar tillbaka ett svar med hur många sekunder klienten har varit igång.
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-boot </b>Om -boot används så får man tillbaka ett timestamp när klienten startades.
    </p>

    <h3>netusage</h3>
    <p>netusage skickar tillbaka klientens nätverks I/O</p>

    <h3>ping [-ports=(port-range)] [-error]</h3>
    <p>ping skickar ett ping meddelande till en eller flera portar för att kolla om dem är öppna, kan användas för att kolla om klienten är igång eller för att kolla vilka portar som används.<br>
    En port range kan ha flera syntax "3333" "22,3333" "22-80,3333".<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-port=(port) </b>-port specifierar vilken eller vilka portar som ping meddelandet ska gå till. <i>Default är 3333.</i><br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-error </b>Om -error används så kommer det att finnas ett error medelande om en eller flera portar misslyckades, annars kommer man bara få error medelande om det är något väldigt dåligt.
    </p>

    <h3>info</h3>
    <p>info skickar tillbaka ett svar med information om klienten så som hostname, OS, interfaces<p>

    <h3>file -file=(file path)</h3>
    <p>file kollar om filen existerar och om den gör det så får man tillbaka information om filen, storlek, senast modifierar, perms.<br>
    <b>&nbsp;&nbsp;&nbsp;&nbsp;-file=(file path) </b>Vart filen ligger som man ska få information från.
    </b>
</div>

<div class="col-md-2"></div>