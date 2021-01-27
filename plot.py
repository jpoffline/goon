import matplotlib.pyplot as plt
import matplotlib.colors as mcolors


def save_to_plot(fff, nx, ny, location="out", prefix="field", timestep=0):
    outfilename = f"{location}/{prefix}_{timestep}.png"
    print(f"saving {outfilename}")
    fig, axs = plt.subplots(2, 2)

    field = []
    for fx in fff:
        ff = []
        for fy in fx:
            ff.append(fy)
        field.append(ff)
    im = axs[0][0].imshow(field)
    fig.colorbar(im, ax=axs[0][0])
    axs[0][0].set_title(f"t = {timestep}")

    sli1 = []
    for i in range(0, nx):
        sli1.append(field[i][int(ny * 0.5)])
    axs[1][0].plot(sli1)

    sli2 = []
    for j in range(0, ny):
        sli2.append(field[int(nx * 0.5)][j])
    axs[0][1].plot(sli2)

    plt.savefig(outfilename)


def read(fn):
    d = []
    with open(fn, "r") as file:
        for l in file.readlines():
            line = l.strip().split(" ")
            r = []
            for i in line:
                r.append(float(i))
            d.append(r)
    return d


t = ["00000", "00020", "00040", "00060", "00180"]
for ts in t:
    fl = read(f"out/field_0_{ts}.dat")
    save_to_plot(fl, 100, 100, prefix="field_0", timestep=ts)
